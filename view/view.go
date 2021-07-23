package view

import (
	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/layout"
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/state"
)

const historySize = 10

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
	Subscribe(topic api.Topic, notify func())
	Unsubscribe()

	SubscribeHandler(handler models.Handler, handle func(string, ...interface{}))
	SubscribeToNamespaces(notify func()) error
	SubscribeToTaskGroups(jobID string, notify func()) error
	SubscribeToLogs(allocID, source string, notify func())

	ForceUpdate()
}

type View struct {
	Client  Client
	Watcher Watcher

	Layout *layout.Layout

	history *History
	state   *state.State

	components *Components

	draw chan struct{}
}

type Components struct {
	ClusterInfo     *component.ClusterInfo
	Selections      *component.Selections
	Commands        *component.Commands
	Logo            *component.Logo
	JobTable        *component.JobTable
	DeploymentTable *component.DeploymentTable
	NamespaceTable  *component.NamespaceTable
	AllocationTable *component.AllocationTable
	TaskGroupTable  *component.TaskGroupTable
	JumpToJob       *component.JumpToJob
	Error           *component.Error
	Info            *component.Info
	Failure         *component.Info
	LogStream       *component.Logger
	LogSearch       *component.SearchField
	Search          *component.SearchField
}

func New(components *Components, watcher Watcher, client Client, state *state.State) *View {
	components.Search = component.NewSearchField()

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
