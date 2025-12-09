// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) Logs(taskName string, allocID, source string) {
	v.Layout.Body.SetTitle(titleLogs)

	v.components.LogSearch.InputField.SetText("")
	v.Layout.Body.Clear()
	v.components.LogStream.Clear()

	logStreamProps := v.components.LogStream.Props
	logStreamProps.TaskName = taskName

	// If the logstream contains data from a previous log stream
	// remove it!
	if len(logStreamProps.Data) > 0 {
		logStreamProps.Data = []byte{}
	}

	v.Layout.Container.SetInputCapture(v.InputLogs)
	v.components.Commands.Update(component.LogCommands)

	update := func() {
		logStreamProps.Data = append(logStreamProps.Data, v.state.Logs...)
		v.components.LogStream.Render()
		v.Draw()
	}

	v.Watcher.SubscribeToLogs(allocID, taskName, source, update)

	v.components.LogStream.ClearDisplay()
	v.components.LogStream.Display()

	update()

	v.Layout.Container.SetFocus(v.components.LogStream.TextView.Primitive())

	v.addToHistory(v.state.SelectedNamespace, models.TopicLog, func() {
		update()
	})

	v.Layout.Container.SetInputCapture(v.InputLogs)
}
