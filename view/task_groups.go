// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) TaskGroups(jobID string) {
	v.Layout.Body.SetTitle(titleTaskGroups)
	v.state.Elements.TableMain = v.components.TaskGroupTable.Table.Primitive().(*tview.Table)

	v.components.Commands.Update(component.NoViewCommands)
	v.Layout.Container.SetInputCapture(v.InputTaskGroups)

	search := v.components.Search

	update := func() {
		v.components.TaskGroupTable.Props.Data = v.state.TaskGroups
		v.components.TaskGroupTable.Props.JobID = jobID
		v.components.TaskGroupTable.Props.HandleNoResources = v.handleNoResources
		v.components.TaskGroupTable.Render()
		v.Draw()
	}

	search.Props.ChangedFunc = func(text string) {
		v.state.Filter.TaskGroups = text
		update()
	}

	v.Watcher.SubscribeToTaskGroups(jobID, update)

	update()

	v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
		v.state.SelectedNamespace = text
		v.TaskGroups(jobID)
	})

	v.addToHistory(v.state.SelectedNamespace, models.TopicTaskGroup, func() {
		v.TaskGroups(jobID)
	})

	v.Layout.Container.SetFocus(v.components.TaskGroupTable.Table.Primitive())
}
