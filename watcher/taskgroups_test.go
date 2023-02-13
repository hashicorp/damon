// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package watcher_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/state"
	"github.com/hcjulz/damon/watcher"
	"github.com/hcjulz/damon/watcher/watcherfakes"
)

func TestSubscribeToTaskGroups_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("It notifies the subscriber initially", func(t *testing.T) {
		// SubscribeToTaskGroups runs a goroutine that polls Nomad based on the given
		// interval to fetch TaskGroups and notify the subscriber.
		// Before the goroutine starts it performs an initial fetch to avoid a delay in the
		// size of the interval duration.
		// In this case we test if this initial call happens.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

		expectedNSFirstCall := []*models.TaskGroup{{Name: "foo"}}

		done := make(chan struct{})

		var callCount int
		notifier := func() {
			callCount++
			switch callCount {
			case 1:
				r.Equal(expectedNSFirstCall, state.TaskGroups)
			case 2:
				// wait for the goroutine to do his first call
				// and finish the test.
				done <- struct{}{}
			}
		}

		nomad.TaskGroupsReturnsOnCall(0, []*models.TaskGroup{
			{Name: "foo"},
		}, nil)

		nomad.TaskGroupsReturnsOnCall(1, []*models.TaskGroup{
			{Name: "foo"},
			{Name: "bar"},
		}, nil)

		watcher.SubscribeToTaskGroups("myJob", notifier)

		<-done
	})

	t.Run("It continues to notify the subscriber after the initial notification", func(t *testing.T) {
		// SubscribeToTaskGroups runs a goroutine that polls Nomad based on the given
		// interval to fetch TaskGroups and notify the subscriber.
		// In this case we test if fetch happes.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

		expectedNSSecondCall := []*models.TaskGroup{{Name: "foo"}, {Name: "bar"}}

		done := make(chan struct{})

		var callCount int
		notifier := func() {
			callCount++
			switch callCount {
			case 3:
				// make sure that the test finishes
				// and avoid blocking if the assertion
				// fails.
				defer func() { done <- struct{}{} }()

				r.Equal(expectedNSSecondCall, state.TaskGroups)
			}
		}

		nomad.TaskGroupsReturnsOnCall(2, []*models.TaskGroup{
			{Name: "foo"},
			{Name: "bar"},
		}, nil)

		watcher.SubscribeToTaskGroups("jobID", notifier)

		<-done

		r.Equal(callCount, 3)
	})

}

func TestSubscribeToTaskGroups_Sad(t *testing.T) {
	// In this case we test that the Error handler
	// is called when nomad returns an error.

	r := require.New(t)

	nomad := &watcherfakes.FakeNomad{}
	state := state.New()
	watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

	var called bool
	watcher.SubscribeHandler(models.HandleError, func(_ string, _ ...interface{}) {
		called = true
	})

	nomad.TaskGroupsReturns(nil, errors.New("argh"))

	watcher.SubscribeToTaskGroups("jobID", func() {})

	r.True(called)
}
