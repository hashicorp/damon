package view

import (
	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) Logs(tasks []string, allocID, source string) {
	v.Layout.Body.SetTitle(titleLogs)

	v.components.LogSearch.InputField.SetText("")
	v.Layout.Body.Clear()
	v.components.LogStream.Clear()

	if len(tasks) > 1 {
		v.components.SelectorModal.Props.Items = tasks
		v.components.SelectorModal.Props.AllocationID = allocID
		v.components.SelectorModal.SetSelectedFunc(func(task string) {
			v.Logs([]string{task}, allocID, source)
		})

		v.components.SelectorModal.Render()

		v.Draw()

		v.Layout.Container.SetFocus(v.components.SelectorModal.Modal.Primitive())

		return
	}

	taskName := tasks[0]

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
