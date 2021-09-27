package view

import (
	"os/exec"

	"github.com/hashicorp/nomad/api"
	"github.com/hcjulz/damon/models"
)

func (v *View) JobStatus(jobID string) {
	v.Layout.Body.Clear()

	v.Layout.Container.SetInputCapture(v.InputMainCommands)
	v.Layout.Container.SetFocus(v.components.JobStatus.TextView.Primitive())
	
	jobStatus := v.components.JobStatus
	jobStatus.Status = "Loading..."
	jobStatus.Render()
	
	update := func() {

		cmd := exec.Command("nomad", "status", jobID)
		stdout,err := cmd.Output()

		if err != nil {
			jobStatus.Status = string(err.Error()) 
		}else {
			jobStatus.Status = string(stdout)
		}

		jobStatus.Render()
		v.Draw()
	}

	v.Watcher.Subscribe(api.TopicJob, update)
	go update()

	v.addToHistory(v.state.SelectedNamespace, models.TopicLog, update)
	v.Layout.Container.SetInputCapture(v.InputMainCommands)
}
