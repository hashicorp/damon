// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"github.com/hcjulz/damon/models"
)

func (v *View) JobStatus(jobID string) {
	v.Layout.Body.SetTitle(titleJobStatus)
	v.Layout.Body.Clear()

	v.Layout.Container.SetInputCapture(v.InputMainCommands)
	v.Layout.Container.SetFocus(v.components.JobStatus.TextView.Primitive())

	jobStatus := v.components.JobStatus

	update := func() {
		jobStatus.Props.Data = v.state.JobStatus

		jobStatus.Render()
		v.Draw()
	}

	v.Watcher.SubscribeToJobStatus(jobID, update)
	update()

	v.addToHistory(v.state.SelectedNamespace, models.TopicLog, update)
	v.Layout.Container.SetInputCapture(v.InputMainCommands)
}
