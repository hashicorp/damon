package view

import (
	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) TaskEvents(allocID string) {
	v.state.Elements.TableMain = v.components.TaskEventsTable.Table.Primitive().(*tview.Table)

	v.components.Commands.Update(component.NoViewCommands)
	v.Layout.Container.SetInputCapture(v.InputMainCommands)

	alloc, ok := v.getAllocation(allocID)
	if !ok {
		v.handleError("allocation with ID %s doesn't exist", allocID)
		return
	}

	tasks := alloc.Tasks
	if len(tasks) == 0 {
		v.handleError("no tasks for allocID %s", allocID)
		return
	}

	// reverse the events array to show latest event on top.
	reverseEvents(tasks[0].Events)

	update := func() {
		v.components.TaskEventsTable.Props.Data = tasks[0].Events
		v.components.TaskEventsTable.Props.AllocID = allocID
		v.components.TaskEventsTable.Props.HandleNoResources = v.handleNoResources
		v.components.TaskEventsTable.Render()
		v.Draw()
	}

	v.Watcher.Subscribe(update, api.TopicAllocation)

	update()

	v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
		v.state.SelectedNamespace = text
		v.TaskEvents(allocID)
	})

	v.addToHistory(v.state.SelectedNamespace, api.TopicAllocation, func() {
		v.TaskEvents(allocID)
	})

	v.Layout.Container.SetFocus(v.components.TaskEventsTable.Table.Primitive())
}

func reverseEvents(e []*api.TaskEvent) {
	for i, j := 0, len(e)-1; i < j; i, j = i+1, j-1 {
		e[i], e[j] = e[j], e[i]
	}
}

func (v *View) getAllocation(id string) (*models.Alloc, bool) {
	for _, a := range v.state.Allocations {
		if a.ID == id {
			return a, true
		}
	}

	return nil, false
}
