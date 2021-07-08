package view

import (
	"regexp"

	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) Allocations(jobID string) {
	v.viewSwitch()

	v.components.Commands.Update(component.AllocCommands)
	v.Layout.Container.SetInputCapture(v.InputAllocations)

	search := v.components.Search
	table := v.components.AllocationTable

	update := func() {
		table.Props.Data = v.filterAllocs(jobID)
		table.Render()
		v.Draw()
	}

	// Overwrite the search change function to filter allocations
	search.Props.ChangedFunc = func(text string) {
		v.state.Filter.Allocations = text
		update()
	}

	v.components.AllocationTable.Props.JobID = jobID

	v.Watcher.Subscribe(api.TopicAllocation, update)

	update()

	v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
		v.state.SelectedNamespace = text
		v.Allocations(jobID)
	})

	// Add this view to the history
	v.addToHistory(v.state.SelectedNamespace, api.TopicAllocation, func() {
		v.Allocations(jobID)
	})

	// Set the current visible table, such that it can be focused when needed
	v.state.Elements.TableMain = v.components.AllocationTable.Table.Primitive().(*tview.Table)

	// focus the current table when loaded
	v.Layout.Container.SetFocus(v.components.AllocationTable.Table.Primitive())
}

func (v *View) filterAllocs(jobID string) []*models.Alloc {
	data := v.filterAllocsForJob(jobID)
	filter := v.state.Filter.Allocations
	if filter != "" {
		rx, _ := regexp.Compile(filter)
		result := []*models.Alloc{}
		for _, alloc := range v.state.Allocations {
			switch true {
			case rx.MatchString(alloc.ID),
				rx.MatchString(alloc.TaskGroup),
				rx.MatchString(alloc.JobID),
				rx.MatchString(alloc.DesiredStatus),
				rx.MatchString(alloc.NodeID),
				rx.MatchString(alloc.NodeName):
				result = append(result, alloc)
			}
		}

		return result
	}

	return data
}

func (v *View) filterAllocsForJob(jobID string) []*models.Alloc {
	rx, _ := regexp.Compile(jobID)
	result := []*models.Alloc{}
	for _, job := range v.state.Allocations {
		switch true {
		case rx.MatchString(job.JobID):
			result = append(result, job)
		}
	}
	return result
}
