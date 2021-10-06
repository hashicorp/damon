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
