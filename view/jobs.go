package view

import (
	"regexp"

	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) Jobs() {
	v.viewSwitch()

	v.Layout.Container.SetInputCapture(v.InputJobs)
	v.components.Commands.Update(component.JobCommands)

	search := v.components.Search
	table := v.components.JobTable

	v.state.Elements.TableMain = table.Table.Primitive().(*tview.Table)

	update := func() {
		table.Props.Data = v.filterJobs()
		table.Props.Namespace = v.state.SelectedNamespace
		table.Render()
		v.Draw()
	}

	search.Props.ChangedFunc = func(text string) {
		v.state.Filter.Jobs = text
		update()
	}

	// TODO
	// go v.Watcher.Watch(func() {
	// 	v.components.JobTable.Render()
	// 	v.Layout.Container.Draw()
	// })

	if table.Props.SelectJob == nil {
		table.Props.SelectJob = func(jobID string) {
			v.Allocations(jobID)
		}
	}

	v.Watcher.Subscribe(api.TopicJob, update)
	if len(v.state.Jobs) == 0 {
		v.Watcher.ForceUpdate()
	}

	update()

	v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
		v.state.SelectedNamespace = text
		v.Jobs()
	})

	v.addToHistory(v.state.SelectedNamespace, api.TopicJob, v.Jobs)
	v.Layout.Container.SetFocus(v.components.JobTable.Table.Primitive())
}

func (v *View) filterJobs() []*models.Job {
	data := v.namespaceFilterJobs()
	filter := v.state.Filter.Jobs
	if filter != "" {
		rx, _ := regexp.Compile(filter)
		result := []*models.Job{}
		for _, job := range data {
			switch true {
			case rx.MatchString(job.ID),
				rx.MatchString(job.Name),
				rx.MatchString(job.Namespace),
				rx.MatchString(job.Status),
				rx.MatchString(job.Type):
				result = append(result, job)
			}
		}

		return result
	}

	return data
}

func (v *View) namespaceFilterJobs() []*models.Job {
	rx, _ := regexp.Compile(v.state.SelectedNamespace)
	result := []*models.Job{}
	for _, job := range v.state.Jobs {
		switch true {
		case rx.MatchString(job.Namespace):
			result = append(result, job)
		}
	}
	return result
}

func (v *View) inputJobs(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyCtrlS:
		jobID := v.components.JobTable.GetIDForSelection()
		v.startStopJob(jobID)
	case tcell.KeyRune:
		switch event.Rune() {
		case 't':
			if v.Layout.Footer.HasFocus() || v.components.Search.InputField.Primitive().HasFocus() {
				return event
			}

			jobID := v.components.JobTable.GetIDForSelection()
			v.TaskGroups(jobID)

		case '/':
			if !v.Layout.Footer.HasFocus() {
				if !v.state.Toggle.Search {
					v.state.Toggle.Search = true
					v.Search()
				} else {
					v.Layout.Container.SetFocus(v.components.Search.InputField.Primitive())
				}
				return nil
			}
		}

	}

	return event
}

func (v *View) startStopJob(jobID string) {
	job, err := v.Client.GetJob(jobID)
	if err != nil {
		v.handleError("failed to start/stop job: %s", err.Error())
		return
	}

	if *job.Status == "dead" {
		v.components.Confirm.Props.Done = func(index int, text string) {
			if index == 1 {
				err := v.Client.StartJob(job)
				v.err(err, "Failed to start job")
			}

			v.closeConfirmModal()
		}

		v.components.Confirm.Render("Do you really want to start the job?")
	} else {
		v.components.Confirm.Props.Done = func(index int, text string) {
			if index == 1 {
				err := v.Client.StopJob(jobID)
				v.err(err, "Failed to stop job")
			}

			v.closeConfirmModal()
		}

		v.components.Confirm.Render("Do you really want to stop the job?")
	}

	v.Layout.Container.SetFocus(v.components.Confirm.Modal.Primitive())
}

func (v *View) closeConfirmModal() {
	v.Layout.Pages.RemovePage(v.components.Confirm.Props.ID)
	v.Layout.Container.SetFocus(v.state.Elements.TableMain)
}
