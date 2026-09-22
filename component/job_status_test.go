// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
	"github.com/hcjulz/damon/models"
)

func TestJobStatus_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is data to render", func(t *testing.T) {
		textView := &componentfakes.FakeTextView{}
		jobStatus := component.NewJobStatus()
		jobStatus.TextView = textView
		jobStatus.Props.Data = &models.JobStatus{
			ID:                "fakeID",
			Name:              "fakeName",
			Namespace:         "fakeNamespace",
			Type:              "fakeType",
			Status:            "fakeStatus",
			StatusDescription: "Some fake description",
			SubmitDate:        time.Date(1987, 10, 16, 12, 45, 0, 0, time.UTC),
			Priority:          0,
			Datacenters:       "fakeDC",
			Periodic:          false,
			Parameterized:     false,
			TaskGroups: []*models.TaskGroup{
				{
					Name:     "fakeGroup",
					JobID:    "fakeJobID",
					Queued:   10,
					Complete: 20,
					Failed:   30,
					Running:  40,
					Starting: 50,
					Lost:     60,
				},
			},
			TaskGroupStatus: []*models.TaskGroupStatus{
				{
					ID:                "fakeGroupStatus",
					Desired:           70,
					Placed:            80,
					Healthy:           90,
					Unhealthy:         100,
					ProgressDeadline:  200,
					Status:            "fake Group Status",
					StatusDescription: "Some fake group description",
				},
			},
			Allocations: []*models.Alloc{
				{
					ID:            "1234567890",
					Name:          "123",
					TaskGroup:     "fakeAllocGroup",
					Tasks:         []models.AllocTask{},
					TaskNames:     []string{},
					JobID:         "",
					JobType:       "",
					NodeID:        "1234567890",
					NodeName:      "",
					DesiredStatus: "fake Desired Status",
					Version:       100,
					Status:        "fake Alloc Status",
					Created:       time.Now(),
					Modified:      time.Now(),
				},
			},
		}

		jobStatus.Bind(tview.NewFlex())

		err := jobStatus.Render()
		r.NoError(err)

		text := textView.SetTextArgsForCall(0)
		r.Equal(strings.ReplaceAll(text, " ", ""), `
ID=fakeID
Name=fakeName
SubmitTime=1987-10-1612:45:00
Type=fakeType
Priority=0
Datacenters=fakeDC
Namespace=fakeNamespace
Status=fakeStatus
Periodic=false
Parameterized=false

Summary
TaskGroupQueuedStartingRunningFailedCompleteLost
fakeGroup105040302060

LatestDeployment
ID=fakeGroupStatus
Status=fakeGroupStatus
StatusDescription=Somefakegroupdescription

Deployed
TaskGroupDesiredPlacedHealthyUnhealthyProgressDeadline
fakeGroupStatus708090100200ns

Allocations
IDNodeIDTaskGroupVersionDesiredStatusCreatedModified
1234567812345678fakeAllocGroup100fakeDesiredStatusfakeAllocStatus0sago0sago
`)
	})

	t.Run("When there is no data to render", func(t *testing.T) {
		textView := &componentfakes.FakeTextView{}
		jobStatus := component.NewJobStatus()
		jobStatus.TextView = textView
		jobStatus.Props.Data = &models.JobStatus{}

		jobStatus.Bind(tview.NewFlex())

		err := jobStatus.Render()
		r.NoError(err)

		text := textView.SetTextArgsForCall(0)
		r.Equal(text, "Status not available.")
	})
}

func TestJobStatus_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the component is not bound", func(t *testing.T) {
		jobStatus := component.NewJobStatus()

		err := jobStatus.Render()
		r.Error(err)

		r.True(errors.Is(err, component.ErrComponentNotBound))
		r.EqualError(err, "component not bound")
	})
}

// TestJobStatus_Tables pins the exact table layout the tablewriter
// configuration produces: borderless, no separators, left aligned, headers
// verbatim and content unwrapped. Trailing spaces are part of the expected
// output, they are the cell padding of the renderer.
func TestJobStatus_Tables(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	jobStatus := component.NewJobStatus()
	jobStatus.TextView = textView

	// No allocations, their columns hold relative timestamps.
	jobStatus.Props.Data = &models.JobStatus{
		ID:          "fakeID",
		Name:        "fakeName",
		Namespace:   "fakeNamespace",
		Type:        "fakeType",
		Status:      "fakeStatus",
		SubmitDate:  time.Date(1987, 10, 16, 12, 45, 0, 0, time.UTC),
		Priority:    50,
		Datacenters: "fakeDC",
		TaskGroups: []*models.TaskGroup{
			{
				Name:     "fakeGroup",
				Queued:   10,
				Starting: 50,
				Running:  40,
				Failed:   30,
				Complete: 20,
				Lost:     60,
			},
		},
		TaskGroupStatus: []*models.TaskGroupStatus{
			{
				ID:                "fakeGroupStatus",
				Desired:           70,
				Placed:            80,
				Healthy:           90,
				Unhealthy:         100,
				ProgressDeadline:  200,
				Status:            "fake Group Status",
				StatusDescription: "Some fake group description",
			},
		},
	}

	jobStatus.Bind(tview.NewFlex())

	err := jobStatus.Render()
	r.NoError(err)

	expected := "\n" +
		" ID             = fakeID              \n" +
		" Name           = fakeName            \n" +
		" SubmitTime     = 1987-10-16 12:45:00 \n" +
		" Type           = fakeType            \n" +
		" Priority       = 50                  \n" +
		" Datacenters    = fakeDC              \n" +
		" Namespace      = fakeNamespace       \n" +
		" Status         = fakeStatus          \n" +
		" Periodic       = false               \n" +
		" Parameterized  = false               \n" +
		"\n  Summary\n" +
		" TaskGroup  Queued  Starting  Running  Failed  Complete  Lost \n" +
		" fakeGroup  10      50        40       30      20        60   \n" +
		"\n  Latest Deployment\n" +
		" ID                  = fakeGroupStatus             \n" +
		" Status              = fake Group Status           \n" +
		" Status Description  = Some fake group description \n" +
		"\n  Deployed\n" +
		" TaskGroup        Desired  Placed  Healthy  Unhealthy  Progress Deadline \n" +
		" fakeGroupStatus  70       80      90       100        200ns             \n" +
		"\n  Allocations\n" +
		" ID  NodeID  TaskGroup  Version  Desired  Status  Created  Modified \n"

	r.Equal(expected, textView.SetTextArgsForCall(0))
}
