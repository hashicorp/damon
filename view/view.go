// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"sync"

	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/layout"
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/state"
)

const (
	historySize = 10

	titleTaskGroups  = "taskgroups"
	titleJobs        = "jobs"
	titleTasks       = "tasks"
	titleJobStatus   = "jobsstatus"
	titleDeployments = "deployments"
	titleNamespaces  = "namespaces"
	titleAllocations = "allocations"
	titleTaskEvents  = "taskevents"
	titleLogs        = "logs"
)

// Client ...
//go:generate counterfeiter . Client
type Client interface {
	GetJob(string) (*api.Job, error)
	StartJob(job *api.Job) error
	StopJob(string) error
}

// Watcher ...
//go:generate counterfeiter . Watcher
type Watcher interface {
	Subscribe(notify func(), topics ...api.Topic)
	Unsubscribe()

	SubscribeHandler(handler models.Handler, handle func(string, ...interface{}))
	SubscribeToNamespaces(notify func())
	SubscribeToTaskGroups(jobID string, notify func()) error
	SubscribeToJobStatus(jobID string, notify func()) error
	SubscribeToLogs(allocID, taskName, source string, notify func())

	ResumeLogs()
}

type View struct {
	Client  Client
	Watcher Watcher

	Layout *layout.Layout

	history *History
	state   *state.State

	components *Components
	mutex      sync.Mutex

	draw chan struct{}
}

type Components struct {
	ClusterInfo     *component.ClusterInfo
	Selections      *component.Selections
	SelectorModal   *component.SelectorModal
	Commands        *component.Commands
	Logo            *component.Logo
	JobTable        *component.JobTable
	TaskTable       *component.TaskTable
	JobStatus       *component.JobStatus
	DeploymentTable *component.DeploymentTable
	NamespaceTable  *component.NamespaceTable
	AllocationTable *component.AllocationTable
	TaskGroupTable  *component.TaskGroupTable
	TaskEventsTable *component.TaskEventsTable
	JumpToJob       *component.JumpToJob
	Error           *component.Error
	Info            *component.Info
	Failure         *component.Info
	LogStream       *component.Logger
	LogSearch       *component.SearchField
	LogHighlight    *component.SearchField
	Search          *component.SearchField
	Confirm         *component.GenericModal
}

func New(components *Components, watcher Watcher, client Client, state *state.State) *View {
	components.Search = component.NewSearchField("")

	return &View{

		history: &History{
			HistorySize: historySize,
		},

		Layout:  layout.New(layout.Default, layout.EnableMouse),
		Watcher: watcher,
		Client:  client,
		state:   state,
		draw:    make(chan struct{}, 1),

		components: components,
	}
}

func (v *View) GoBack() {
	v.history.pop()
}

func (v *View) JumpToJob() {
	jump := v.components.JumpToJob
	jump.Props.Jobs = v.state.Jobs
	v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 1)
	jump.Render()
	v.Layout.Container.SetFocus(jump.InputField.Primitive())
}

func (v *View) LogSearch() {
	search := v.components.LogSearch
	v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 1)
	search.Render()
	v.Layout.Container.SetFocus(search.InputField.Primitive())
}

func (v *View) LogHighlight() {
	highlight := v.components.LogHighlight
	v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 1)
	highlight.Render()
	v.Layout.Container.SetFocus(highlight.InputField.Primitive())
}

func (v *View) Search() {
	search := v.components.Search
	v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 1)
	search.Render()
	v.Layout.Container.SetFocus(search.InputField.Primitive())
}

func (v *View) Draw() {
	v.draw <- struct{}{}
}

// DrawLoop refreshes the screen when it receives a
// signal on the draw channel. This function should
// be run inside a goroutine as tview.Application.Draw()
// can deadlock when called from the main thread.
func (v *View) DrawLoop(stop chan struct{}) {
	for {
		select {
		case <-v.draw:
			v.Layout.Container.Draw()
		case <-stop:
			return
		}
	}
}

func (v *View) addToHistory(ns string, topic api.Topic, update func()) {
	v.history.push(func() {
		v.state.SelectedNamespace = ns
		// update()

		// v.components.Selections.Props.Rerender = update
		v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
			v.state.SelectedNamespace = text
			update()
		})
		// v.Watcher.Subscribe(topic, update)

		index := getNamespaceNameIndex(ns, v.state.Namespaces)
		v.state.Elements.DropDownNamespace.SetCurrentOption(index)
	})
}

func (v *View) viewSwitch() {
	v.resetSearch()
}
