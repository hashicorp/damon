// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
	"github.com/hcjulz/damon/models"
)

// renderer is implemented by every component that renders into a table slot.
type renderer interface {
	Bind(slot *tview.Flex)
	Render() error
}

// idWithVerbs contains printf verbs, Table.SetTitle formats its arguments,
// so an ID like this must reach the title unchanged.
const idWithVerbs = "job-100%s-%d"

func TestTableTitles(t *testing.T) {
	noResources := func(format string, args ...interface{}) {}

	tests := []struct {
		name  string
		title string
		build func(table component.Table) renderer
	}{
		{
			name:  "allocations",
			title: fmt.Sprintf("%s (Job: %s)", component.TableTitleAllocations, idWithVerbs),
			build: func(table component.Table) renderer {
				c := component.NewAllocationTable()
				c.Table = table
				c.Props.SelectAllocation = func(allocID string) {}
				c.Props.HandleNoResources = noResources
				c.Props.JobID = idWithVerbs
				c.Props.Data = []*models.Alloc{{ID: "alloc-id"}}
				return c
			},
		},
		{
			name:  "deployments",
			title: fmt.Sprintf("%s (%s)", component.TableTitleDeployments, idWithVerbs),
			build: func(table component.Table) renderer {
				c := component.NewDeploymentTable()
				c.Table = table
				c.Props.SelectDeployment = func(id string) {}
				c.Props.HandleNoResources = noResources
				c.Props.Namespace = idWithVerbs
				c.Props.Data = []*models.Deployment{{ID: "deployment-id"}}
				return c
			},
		},
		{
			name:  "task events",
			title: fmt.Sprintf("%s (%s)", component.TableTitleTaskEvents, idWithVerbs),
			build: func(table component.Table) renderer {
				c := component.NewTaskEventsTable()
				c.Table = table
				c.Props.HandleNoResources = noResources
				c.Props.AllocID = idWithVerbs
				c.Props.Data = []*api.TaskEvent{{Type: "type"}}
				return c
			},
		},
		{
			name:  "task groups",
			title: fmt.Sprintf("%s (%s)", component.TableTitleTaskGroups, idWithVerbs),
			build: func(table component.Table) renderer {
				c := component.NewTaskGroupTable()
				c.Table = table
				c.Props.SelectTaskGroup = func(id string) {}
				c.Props.HandleNoResources = noResources
				c.Props.JobID = idWithVerbs
				c.Props.Data = []*models.TaskGroup{{Name: "group"}}
				return c
			},
		},
		{
			name:  "tasks",
			title: fmt.Sprintf("%s (Allocation: %s)", component.TableTitleTasks, idWithVerbs),
			build: func(table component.Table) renderer {
				c := component.NewTaskTable()
				c.Table = table
				c.Props.SelectTask = func(allocID, taskID string) {}
				c.Props.HandleNoResources = noResources
				c.Props.AllocationID = idWithVerbs
				c.Props.Data = []*models.Task{{
					Name:   "task",
					Events: []*api.TaskEvent{{DisplayMessage: "msg"}},
				}}
				return c
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := require.New(t)

			fakeTable := &componentfakes.FakeTable{}

			c := test.build(fakeTable)
			c.Bind(tview.NewFlex())

			err := c.Render()
			r.NoError(err)

			r.Equal(1, fakeTable.SetTitleCallCount())

			format, args := fakeTable.SetTitleArgsForCall(0)
			r.Equal(test.title, fmt.Sprintf(format, args...))
		})
	}
}
