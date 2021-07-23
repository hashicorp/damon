package view

import (
	"regexp"

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

func (v *View) viewSwitch() {
	v.resetSearch()
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
