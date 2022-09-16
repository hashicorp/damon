package view

import (
	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) Tasks(alloc *models.Alloc) {
	v.viewSwitch()
	v.Layout.Body.SetTitle(titleTasks)

	v.Layout.Container.SetInputCapture(v.InputMainCommands)
	v.components.Commands.Update(component.TaskCommands)

	table := v.components.TaskTable
	table.Props.Data = alloc.TaskList
	table.Props.AllocationID = alloc.ID

	v.state.Elements.TableMain = table.Table.Primitive().(*tview.Table)

	update := func() {
		table.Render()

		v.Draw()
	}

	if table.Props.SelectTask == nil {
		table.Props.SelectTask = func(taskName, allocID string) {
			v.Logs(taskName, allocID, "stdout")
		}
	}

	v.Watcher.Subscribe(update, api.TopicAllocation)

	update()

	v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
		v.state.SelectedNamespace = text
		v.Tasks(alloc)
	})

	v.addToHistory(v.state.SelectedNamespace, api.TopicAllocation, func() {
		v.Tasks(alloc)
	})

	v.Layout.Container.SetFocus(v.components.TaskTable.Table.Primitive())
}
