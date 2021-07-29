package watcher

import (
	"time"

	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/state"
)

//go:generate counterfeiter . Activities
type Activities interface {
	Add(chan struct{})
	DeactivateAll()
}

//go:generate counterfeiter . Nomad
type Nomad interface {
	Address() string
	Jobs(*nomad.SearchOptions) ([]*models.Job, error)
	Namespaces(*nomad.SearchOptions) ([]*models.Namespace, error)
	Deployments(*nomad.SearchOptions) ([]*models.Deployment, error)
	TaskGroups(string, *nomad.SearchOptions) ([]*models.TaskGroup, error)
	Allocations(*nomad.SearchOptions) ([]*models.Alloc, error)
	JobAllocs(string, *nomad.SearchOptions) ([]*models.Alloc, error)
	Logs(allocID, taskNmae, logType string, cancel <-chan struct{}) (<-chan *api.StreamFrame, <-chan error)
	Stream(topics nomad.Topics, index uint64) (<-chan *api.Events, error)
}

// Watcher watches a Nomad cluster for updates and
// updates the central state accordingly. Whenever
// an update happens it notifies the current subscriber.
type Watcher struct {
	state      *state.State
	subscriber *subscriber
	handlers   map[models.Handler]func(msg string, args ...interface{})
	nomad      Nomad

	forceUpdate chan api.Topic
	activities  Activities

	interval time.Duration
}

type subscriber struct {
	topic  api.Topic
	notify func()
}

func NewWatcher(state *state.State, nomad Nomad, interval time.Duration) *Watcher {
	return &Watcher{
		state:       state,
		nomad:       nomad,
		handlers:    map[models.Handler]func(ms string, args ...interface{}){},
		forceUpdate: make(chan api.Topic),
		activities:  &ActivityPool{},
		interval:    interval,
	}
}

// Subscribe subscribes a function to a topic. This function should always be
// called before Watcher.activities.Add().
func (w *Watcher) Subscribe(topic api.Topic, notify func()) {
	w.subscriber = &subscriber{
		topic:  topic,
		notify: notify,
	}

	// Whenever a subscription comes in make sure all running
	// goroutines (expect the main (Watch)) are stopped.
	// This is because currently Damon uses event streams for
	// Deployments and Jobs, but polls the API for Namespaces,
	// Allocations, and TaskGroups.
	w.activities.DeactivateAll()
}

// Unsubscribe removes the current subscriber.
func (w *Watcher) Unsubscribe() {
	w.subscriber = nil
}

// SubscribeHandler subscribes a handler to the watcher. This can be an for example an error
// handler. The handler types are defined in the models package.
func (w *Watcher) SubscribeHandler(handler models.Handler, handle func(string, ...interface{})) {
	w.handlers[handler] = handle
}

// NotifyHandler notifies a handler that an event occurred
// on the topic it subscribed for.
func (w *Watcher) NotifyHandler(handler models.Handler, msg string, args ...interface{}) {
	if _, ok := w.handlers[handler]; ok {
		w.handlers[handler](msg, args...)
	}
}

// Notify notifies the current subscriber on a specific topic (eg Jobs)
// that data got updated in the state.
func (w *Watcher) Notify(topic api.Topic) {
	if w.subscriber != nil && w.subscriber.notify != nil {
		if w.subscriber.topic == topic {
			w.subscriber.notify()
		}
	}
}

// Watch starts an Nomad event stream for top level objects,
// such as Jobs and Deployments. Subscribers to a specific topic
// get notified when an event occurs.
func (w *Watcher) Watch() {
	topics := map[api.Topic][]string{
		api.TopicJob:        {"*"},
		api.TopicDeployment: {"*"},
		api.TopicAllocation: {"*"},
	}

	w.update(api.TopicJob)
	w.update(api.TopicDeployment)
	w.update(api.TopicAllocation)

	index := uint64(1000)
	eventCh, err := w.nomad.Stream(topics, index)
	if err != nil {
		w.NotifyHandler(models.HandleFatal, err.Error())
	}

	for {
		select {
		case event, open := <-eventCh:
			if !open {
				w.NotifyHandler(models.HandleFatal, "event stream closed")
				return
			}

			for _, e := range event.Events {
				w.update(e.Topic)
			}
		case topic := <-w.forceUpdate:
			w.update(topic)

		}
	}
}

// ForceUpdate forces the event stream loop
// to update the state and notify the current
// subscriber.
func (w *Watcher) ForceUpdate() {
	w.forceUpdate <- w.subscriber.topic
}

func (w *Watcher) update(topic api.Topic) {
	switch topic {
	case api.TopicJob:
		w.updateJobs()
	case api.TopicAllocation:
		w.updateAllocations()
	case api.TopicDeployment:
		w.updateDeployments()
	}

	w.Notify(topic)
}

func (w *Watcher) updateJobs() {
	jobs, err := w.nomad.Jobs(&nomad.SearchOptions{
		Namespace: "*",
	})
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
	}

	w.state.Jobs = jobs
}

func (w *Watcher) updateDeployments() {
	dep, err := w.nomad.Deployments(&nomad.SearchOptions{})
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
	}

	w.state.Deployments = dep
}

func (w *Watcher) updateAllocations() {
	allocs, err := w.nomad.Allocations(&nomad.SearchOptions{Namespace: "*"})
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
	}

	w.state.Allocations = allocs
}
