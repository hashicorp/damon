package watcher_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/state"
	"github.com/hcjulz/damon/watcher"
	"github.com/hcjulz/damon/watcher/watcherfakes"
)

func TestSubscribeToLogs_Happy(t *testing.T) {
	// In this case we test the happy path of the
	// function. This includes the initial notification
	// and checking the stream is called for the correct
	// allocation and task.

	r := require.New(t)

	nomad := &watcherfakes.FakeNomad{}
	state := state.New()
	state.Allocations = []*models.Alloc{
		{
			ID:        "the-alloc",
			TaskNames: []string{"the-task"},
		},
		{
			ID:        "another-alloc",
			TaskNames: []string{"another-task"},
		},
	}
	state.Logs = []byte("an initial log line that should be wiped")

	watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

	// A new subscription will close the cancel chan
	defer watcher.Subscribe(func() {}, api.TopicJob)

	streamChan := make(chan *api.StreamFrame)
	errChan := make(chan error)

	nomad.LogsReturns(streamChan, errChan)

	var callCount int
	notify := func() {
		callCount++
	}

	watcher.SubscribeToLogs("the-alloc", "the-task", "stderr", notify)

	actualAllocID, taskName, actualSource, _ := nomad.LogsArgsForCall(0)

	// Check that the initial call happened for the right task and allocation.
	r.Equal("the-alloc", actualAllocID)
	r.Equal("the-task", taskName)
	r.Equal("stderr", actualSource)
	r.Equal(callCount, 1)
	r.Nil(state.Logs)

	// Next the goroutine should notify the subscriber whenever a new logs is available.
	streamChan <- &api.StreamFrame{
		Data: []byte("a new log line\n"),
	}

	r.Eventually(func() bool {
		return callCount == 2
	}, time.Second*5, time.Microsecond*5)

	r.Equal(state.Logs, []byte("a new log line\n"))

	// further log lines should be appended
	streamChan <- &api.StreamFrame{
		Data: []byte("another log line\n"),
	}

	r.Eventually(func() bool {
		return callCount == 3
	}, time.Second*5, time.Microsecond*5)

	r.Equal(state.Logs, []byte("a new log line\nanother log line\n"))
}

func TestSubscribeToLogs_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("the allocation does not exist", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		state.Allocations = []*models.Alloc{}

		watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

		var called bool
		watcher.SubscribeHandler(models.HandleError, func(msg string, args ...interface{}) {
			r.Equal(fmt.Sprintf(msg, args...), "allocation not found: alloc-id")
			called = true
		})

		var callCount int
		watcher.SubscribeToLogs("alloc-id", "task-id", "some-source", func() { callCount++ })

		r.True(called)
		r.Equal(callCount, 0)
	})

	t.Run("When the allocation does not contain any task names", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		state.Allocations = []*models.Alloc{
			{ID: "alloc-id"},
		}

		watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

		var called bool
		watcher.SubscribeHandler(models.HandleError, func(msg string, args ...interface{}) {
			r.Equal(fmt.Sprintf(msg, args...), "no tasks for allocation: alloc-id")
			called = true
		})

		var callCount int
		watcher.SubscribeToLogs("alloc-id", "task-id", "some-source", func() { callCount++ })

		r.True(called)
		r.Equal(callCount, 0)
	})

	t.Run("the error channel receives an error", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		state.Allocations = []*models.Alloc{
			{ID: "alloc-id", TaskNames: []string{"task"}},
		}

		watcher := watcher.NewWatcher(state, nomad, time.Millisecond*250)

		defer watcher.Subscribe(func() {}, api.TopicJob)

		streamChan := make(chan *api.StreamFrame)
		errChan := make(chan error)

		nomad.LogsReturns(streamChan, errChan)

		var called bool
		watcher.SubscribeHandler(models.HandleError, func(msg string, args ...interface{}) {
			r.Equal(fmt.Sprintf(msg, args...), "streaming error")
			called = true
		})

		var callCount int
		watcher.SubscribeToLogs("alloc-id", "task-id", "some-source", func() { callCount++ })

		r.Equal(callCount, 1)

		errChan <- errors.New("streaming error")

		r.Eventually(func() bool {
			return called
		}, time.Second*5, time.Microsecond*5)
	})
}
