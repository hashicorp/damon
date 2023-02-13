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

func TestSubscribeToJobStatus_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("It notifies the subscriber initially", func(t *testing.T) {
		// SubscribeToJobStatus runs a goroutine that polls Nomad based on the given
		// interval to fetch JobStatus and notify the subscriber (if subscribed to JobStatus).
		// Before the goroutine starts it performs an initial fetch to avoid a delay in the
		// size of the interval duration.
		// In this case we test if this initial call happens.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

		expectedNSFirstCall := &models.JobStatus{Name: "foo"}

		done := make(chan struct{})

		var callCount int
		notifier := func() {
			callCount++
			switch callCount {
			case 1:
				r.Equal(expectedNSFirstCall, state.JobStatus)
			case 2:
				// wait for the goroutine to do his first call
				// and finish the test.
				done <- struct{}{}
			}
		}

		nomad.JobStatusReturnsOnCall(0, &models.JobStatus{
			Name: "foo",
		}, nil)

		nomad.JobStatusReturnsOnCall(1, &models.JobStatus{
			Name: "bar",
		}, nil)

		watcher.SubscribeToJobStatus("myJob", notifier)

		<-done
	})

	t.Run("It continues to notify the subscriber after the initial notification", func(t *testing.T) {
		// SubscribeToJobStatus runs a goroutine that polls Nomad based on the given
		// interval to fetch JobStatus and notify the subscriber (if subscribed to JobStatus).
		// In this case we test if fetch happes.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

		expectedNSSecondCall := &models.JobStatus{Name: "foo"}

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

				r.Equal(expectedNSSecondCall, state.JobStatus)
			}
		}

		nomad.JobStatusReturnsOnCall(2, &models.JobStatus{
			Name: "foo",
		}, nil)

		watcher.SubscribeToJobStatus("myJob", notifier)

		<-done

		r.Equal(callCount, 3)
	})

}

func TestSubscribeToJobStatus_Sad(t *testing.T) {
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

	nomad.JobStatusReturns(nil, errors.New("argh"))

	watcher.SubscribeToJobStatus("myJob", func() {})

	r.True(called)
}
