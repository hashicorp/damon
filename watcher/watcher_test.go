// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package watcher_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/state"
	"github.com/hcjulz/damon/watcher"
	"github.com/hcjulz/damon/watcher/watcherfakes"
)

func TestSubscription(t *testing.T) {
	r := require.New(t)

	nomad := &watcherfakes.FakeNomad{}
	state := state.New()

	watcher := watcher.NewWatcher(state, nomad, time.Second*2)

	var called bool
	fn := func() {
		called = true
	}

	watcher.Subscribe(fn, api.TopicJob)
	watcher.Notify(api.TopicJob)

	r.True(called)

	called = false
	watcher.Unsubscribe()
	watcher.Notify(api.TopicJob)

	r.False(called)
}

func TestHandlerSubscription(t *testing.T) {
	r := require.New(t)

	nomad := &watcherfakes.FakeNomad{}
	state := state.New()

	watcher := watcher.NewWatcher(state, nomad, time.Second*2)

	var calledErrHandler bool
	handleErr := func(_ string, _ ...interface{}) {
		calledErrHandler = true
	}

	var calledInfoHandler bool
	handleInfo := func(_ string, _ ...interface{}) {
		calledInfoHandler = true
	}

	var calledFatalHandler bool
	handleFatal := func(_ string, _ ...interface{}) {
		calledFatalHandler = true
	}

	watcher.SubscribeHandler(models.HandleError, handleErr)
	watcher.SubscribeHandler(models.HandleInfo, handleInfo)
	watcher.SubscribeHandler(models.HandleFatal, handleFatal)

	watcher.NotifyHandler(models.HandleError, "error")
	watcher.NotifyHandler(models.HandleInfo, "info")
	watcher.NotifyHandler(models.HandleFatal, "fatal")

	r.True(calledErrHandler)
	r.True(calledInfoHandler)
	r.True(calledFatalHandler)
}

func TestWatch_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is one subscriber for a topic", func(t *testing.T) {
		// In this case, we expect that a subscriber that is subscribed
		// to one of the following topics: Job, Deployment, Allocation
		// is called whenever a event comes in.
		// We expect that the subscriber is notfied initially before the
		// stream starts and the state gets updated accordingly.

		// Setup the watcher
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		// handleErr := func(_ string, _ ...interface{}) {}
		// watcher.SubscribeHandler(models.HandleError, handleErr)

		// Setup expectations
		expectedJobsInitialCall := []*models.Job{{ID: "jupiter"}}
		expectedJobsUpdated := []*models.Job{{ID: "jupiter"}, {ID: "saturn"}}

		// callCount indicates how often the subscriber was notified
		var callCount int
		notifier := func() {
			callCount++
		}

		// Subscribe the notifier to the Job topic
		watcher.Subscribe(notifier, api.TopicJob)

		// Create an eventCh we can send events to for testing...
		eventCh := make(chan *api.Events)
		defer close(eventCh)

		// ...and let the fake nomad client return it.
		nomad.StreamReturns(eventCh, nil)

		// Declare what the the fake client should return on the different calls
		nomad.JobsReturnsOnCall(0, []*models.Job{{ID: "jupiter"}}, nil)
		nomad.JobsReturnsOnCall(1, []*models.Job{
			{ID: "jupiter"},
			{ID: "saturn"},
		}, nil)

		go watcher.Watch()

		r.Eventually(func() bool {
			return callCount == 1
		}, time.Second*5, time.Microsecond*5)

		r.Equal(expectedJobsInitialCall, state.Jobs)

		// We send events for all three topics we are intrested in:
		events := &api.Events{
			Events: []api.Event{
				{
					Topic: api.TopicJob,
				},
			},
		}
		eventCh <- events

		// We expected that the callCount eventually was called twice:
		r.Eventually(func() bool {
			return callCount == 2
		}, time.Second*5, time.Microsecond*5)

		r.Equal(expectedJobsUpdated, state.Jobs)

		// Check that the call counts for each function haven't been called
		// more often than expected.
		r.Equal(nomad.JobsCallCount(), 2)
	})

	t.Run("When the subscriber subscribes for multiple topics", func(t *testing.T) {
		// In this case, we expect that a subscriber that is subscribed
		// to multiple topics: Job & Allocation
		// is called for all events that arrive.
		// We expect that the subscriber is notfied initially for both topics
		// before the stream starts and the state gets updated accordingly.

		// Setup the watcher
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		// Setup expectations
		expectedJobsInitialCall := []*models.Job{{ID: "jupiter"}}
		expectedJobsUpdated := []*models.Job{{ID: "jupiter"}, {ID: "saturn"}}
		expectedAllocsInitialCall := []*models.Alloc{{ID: "bumblebee"}}
		expectedAllocsUpdated := []*models.Alloc{{ID: "prime"}, {ID: "megatron"}}

		// callCount indicates how often the subscriber was notified
		var callCount int
		notifier := func() {
			callCount++
		}

		// Subscribe the notifier to the Job topic and the Allocation topic
		watcher.Subscribe(notifier, api.TopicJob, api.TopicAllocation)

		// Create an eventCh we can send events to for testing...
		eventCh := make(chan *api.Events)
		defer close(eventCh)

		// ...and let the fake nomad client return it.
		nomad.StreamReturns(eventCh, nil)

		// Declare what the the fake client should return on the different calls
		nomad.JobsReturnsOnCall(0, []*models.Job{{ID: "jupiter"}}, nil)
		nomad.JobsReturnsOnCall(1, []*models.Job{
			{ID: "jupiter"},
			{ID: "saturn"},
		}, nil)

		nomad.AllocationsReturnsOnCall(0, []*models.Alloc{{ID: "bumblebee"}}, nil)
		nomad.AllocationsReturnsOnCall(1, []*models.Alloc{
			{ID: "prime"},
			{ID: "megatron"},
		}, nil)

		go watcher.Watch()

		r.Eventually(func() bool {
			return callCount == 2
		}, time.Second*5, time.Microsecond*5)

		r.Equal(expectedJobsInitialCall, state.Jobs)
		r.Equal(expectedAllocsInitialCall, state.Allocations)

		// We send events for all three topics we are intrested in:
		events := &api.Events{
			Events: []api.Event{
				{
					Topic: api.TopicJob,
				},
				{
					Topic: api.TopicAllocation,
				},
			},
		}
		eventCh <- events

		// We expected that the callCount eventually was called twice:
		r.Eventually(func() bool {
			return callCount == 4
		}, time.Second*5, time.Microsecond*5)

		r.Equal(expectedJobsUpdated, state.Jobs)
		r.Equal(expectedAllocsUpdated, state.Allocations)

		r.Equal(nomad.JobsCallCount(), 2)
		r.Equal(nomad.AllocationsCallCount(), 2)
	})

	t.Run("When a subscriber subscribes to a specific topic it doesn't get notified for other topics", func(t *testing.T) {
		// In this case we expect that a subscriber only gets notified for
		// the topic it is subscribed to. Other topics shouldn't call the subscriber.

		// Setup the watcher
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		// callCount indicates how often the subscriber was notified
		var callCount int
		notifier := func() {
			callCount++
		}

		// Subscribe the notifier to the Deployment topic
		watcher.Subscribe(notifier, api.TopicDeployment)

		// Create an eventCh we can send events to for testing...
		eventCh := make(chan *api.Events)
		defer close(eventCh)

		// ...and let the fake nomad client return it.
		nomad.StreamReturns(eventCh, nil)

		expectedJobs := []*models.Job{{ID: "jupiter"}, {ID: "saturn"}}
		nomad.JobsReturnsOnCall(1, expectedJobs, nil)

		go watcher.Watch()

		// We send events for te Job topic.
		events := &api.Events{
			Events: []api.Event{
				{
					Topic: api.TopicJob,
				},
			},
		}
		eventCh <- events

		r.Eventually(func() bool {
			return len(state.Jobs) == len(expectedJobs)
		}, time.Second*5, time.Microsecond*5)

		// We make sure deployments did not get notified.
		r.Equal(callCount, 1)
		r.Equal(expectedJobs, state.Jobs)
	})

	t.Run("When events come in for different topics the state for all topcis gets updated", func(t *testing.T) {
		// In this case we expect when events arrive for Job, Deployment, and Allocation the
		// the corresponding state gets updated.

		// Setup the watcher
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		// Create an eventCh we can send events to for testing...
		eventCh := make(chan *api.Events)
		defer close(eventCh)

		// ...and let the fake nomad client return it.
		nomad.StreamReturns(eventCh, nil)

		// Declare what the the fake client should return on the different calls
		nomad.JobsReturnsOnCall(1, []*models.Job{
			{ID: "jupiter"},
			{ID: "saturn"},
		}, nil)

		nomad.DeploymentsReturnsOnCall(1, []*models.Deployment{
			{ID: "jupiter"},
			{ID: "saturn"},
		}, nil)

		nomad.AllocationsReturnsOnCall(1, []*models.Alloc{
			{ID: "jupiter"},
			{ID: "saturn"},
		}, nil)

		go watcher.Watch()

		// We send events for all three topics we are intrested in:
		events := &api.Events{
			Events: []api.Event{
				{
					Topic: api.TopicJob,
				},
				{
					Topic: api.TopicDeployment,
				},
				{
					Topic: api.TopicAllocation,
				},
			},
		}

		eventCh <- events

		// Wait till the updates happened
		r.Eventually(func() bool {
			return nomad.AllocationsCallCount() == 2 &&
				nomad.DeploymentsCallCount() == 2 &&
				nomad.JobsCallCount() == 2
		}, time.Second*5, time.Microsecond*5)

		expectedJobs := []*models.Job{{ID: "jupiter"}, {ID: "saturn"}}
		expectedDeployments := []*models.Deployment{{ID: "jupiter"}, {ID: "saturn"}}
		expectedAllocs := []*models.Alloc{{ID: "jupiter"}, {ID: "saturn"}}

		r.Equal(expectedJobs, state.Jobs)
		r.Equal(expectedDeployments, state.Deployments)
		r.Equal(expectedAllocs, state.Allocations)
	})

	t.Run("When the subscriber switches", func(t *testing.T) {
		// In this case we expect that the previous
		// subscriber doesn't get notified anymore.
		// Only the new subscriber gets notified
		// for the topic it subscribed to.

		// Setup the watcher
		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		// callCount indicates how often the subscriber was notified
		var callCountJob int
		notifierJob := func() {
			callCountJob++
		}

		var callCountDepl int
		notifierDepl := func() {
			callCountDepl++
		}

		// Create an eventCh we can send events to for testing...
		eventCh := make(chan *api.Events)
		defer close(eventCh)

		// ...and let the fake nomad client return it.
		nomad.StreamReturns(eventCh, nil)

		watcher.Subscribe(notifierDepl, api.TopicDeployment)

		go watcher.Watch()

		// We send events for te both topics.
		events := &api.Events{
			Events: []api.Event{
				{Topic: api.TopicJob},
				{Topic: api.TopicDeployment},
			},
		}

		eventCh <- events

		r.Eventually(func() bool {
			// We expect the deployment callcount is 2
			// because of the initical update call
			// before the even stream starts.
			return callCountDepl == 2
		}, time.Second*5, time.Microsecond*5)

		r.Equal(callCountJob, 0)

		// We overwrite the subscriber
		watcher.Subscribe(notifierJob, api.TopicJob)

		// Send events again
		eventCh <- events

		r.Eventually(func() bool {
			return callCountJob == 1
		}, time.Second*5, time.Microsecond*5)

		// Both notifier should have been called max once.
		r.Equal(callCountJob, 1)

		// The Deployment Call Count should remain the same.
		r.Equal(callCountDepl, 2)
	})
}

func TestWatch_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the event stream can't be setup", func(t *testing.T) {
		// In this case Damon should call the Fatal handler.
		// Which means Damon will be terminated.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		nomad.StreamReturns(nil, errors.New("argh"))

		var called bool
		var message string
		handleFatal := func(msg string, _ ...interface{}) {
			called = true
			message = msg
		}

		watcher.SubscribeHandler(models.HandleFatal, handleFatal)

		go watcher.Watch()

		r.Eventually(func() bool {
			return called
		}, time.Second*5, time.Microsecond*5)

		r.Equal(message, "argh")
	})

	t.Run("When the event stream got closed", func(t *testing.T) {
		// In this case Damon should call the Fatal handler.
		// Which means Damon will be terminated.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		eventCh := make(chan *api.Events)

		nomad.StreamReturns(eventCh, nil)

		var called bool
		var message string
		handleFatal := func(msg string, _ ...interface{}) {
			called = true
			message = msg
		}

		watcher.SubscribeHandler(models.HandleFatal, handleFatal)

		go watcher.Watch()

		// close the eventCh to cause an error
		close(eventCh)

		r.Eventually(func() bool {
			return called
		}, time.Second*5, time.Microsecond*5)

		r.Equal(message, "event stream closed")
	})

	t.Run("When jobs can't be fetched from the nomad cluster", func(t *testing.T) {
		// In this case Damon should call the Error handler.
		// Which means Damon will not be terminated.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		eventCh := make(chan *api.Events)
		defer close(eventCh)

		nomad.StreamReturns(eventCh, nil)

		var called bool
		var message string
		handleErr := func(msg string, _ ...interface{}) {
			called = true
			message = msg
		}

		watcher.SubscribeHandler(models.HandleError, handleErr)

		nomad.JobsReturns(nil, errors.New("argh"))

		go watcher.Watch()

		r.Eventually(func() bool {
			return called
		}, time.Second*5, time.Microsecond*5)

		r.Equal(message, "argh")
	})

	t.Run("When deployments can't be fetched from the nomad cluster", func(t *testing.T) {
		// In this case Damon should call the Error handler.
		// Which means Damon will not be terminated.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		eventCh := make(chan *api.Events)
		defer close(eventCh)

		nomad.StreamReturns(eventCh, nil)

		var called bool
		var message string
		handleErr := func(msg string, _ ...interface{}) {
			called = true
			message = msg
		}

		watcher.SubscribeHandler(models.HandleError, handleErr)

		nomad.DeploymentsReturns(nil, errors.New("argh"))

		go watcher.Watch()

		r.Eventually(func() bool {
			return called
		}, time.Second*5, time.Microsecond*5)

		r.Equal(message, "argh")
	})

	t.Run("When allocations can't be fetched from the nomad cluster", func(t *testing.T) {
		// In this case Damon should call the Error handler.
		// Which means Damon will not be terminated.

		nomad := &watcherfakes.FakeNomad{}
		state := state.New()
		watcher := watcher.NewWatcher(state, nomad, time.Second*2)

		eventCh := make(chan *api.Events)
		defer close(eventCh)

		nomad.StreamReturns(eventCh, nil)

		var called bool
		var message string
		handleErr := func(msg string, _ ...interface{}) {
			called = true
			message = msg
		}

		watcher.SubscribeHandler(models.HandleError, handleErr)

		nomad.AllocationsReturns(nil, errors.New("argh"))

		go watcher.Watch()

		r.Eventually(func() bool {
			return called
		}, time.Second*5, time.Microsecond*5)

		r.Equal(message, "argh")
	})
}
