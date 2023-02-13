// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
	"github.com/stretchr/testify/require"
)

func TestJobStatus(t *testing.T) {
	r := require.New(t)

	fakeJobClient := &nomadfakes.FakeJobClient{}
	fakeDpClient := &nomadfakes.FakeDeploymentClient{}
	client := &nomad.Nomad{JobClient: fakeJobClient, DpClient: fakeDpClient}

	t.Run("When there are no issues", func(t *testing.T) {
		now := time.Now().UnixNano()
		nowUnix := time.Unix(0, now)

		id := "fakeID"
		name := "fakeName"
		fakeType := "fakeType"
		status := "fakestatus"
		value := 100
		ns := "ns"
		fakeTaskGroup := "fakeTaskGroup"

		fakeJobClient.InfoReturns(&api.Job{
			ID:         &id,
			Name:       &name,
			Type:       &fakeType,
			Priority:   &value,
			Status:     &status,
			SubmitTime: &now,
			Namespace:  &ns,
			TaskGroups: []*api.TaskGroup{
				{
					Name: &fakeTaskGroup,
				},
			},
		}, nil, nil)

		fakeJobClient.SummaryReturns(&api.JobSummary{
			JobID: "fakejobId",
			Summary: map[string]api.TaskGroupSummary{
				"Mandalorian": {
					Running: 1,
					Queued:  1,
				},
				"Grogu": {
					Failed: 1,
					Queued: 1,
				},
			},
		}, nil, nil)

		fakeDpClient.ListReturns([]*api.Deployment{
			{
				JobID: "fakeID",
				TaskGroups: map[string]*api.DeploymentState{
					"dp1": {
						HealthyAllocs:   1,
						UnhealthyAllocs: 0,
					},
				},
			},
		}, nil, nil)

		so := nomad.SearchOptions{Namespace: "default"}

		jobStatus, err := client.JobStatus("fakeID", &so)

		jId, queryOptions := fakeJobClient.InfoArgsForCall(0)

		expectedJobStatus := &models.JobStatus{
			ID:         id,
			Name:       name,
			Type:       fakeType,
			Status:     status,
			SubmitDate: nowUnix,
			Namespace:  ns,
			Priority:   value,
			TaskGroups: []*models.TaskGroup{
				{
					Name:     "Grogu",
					JobID:    "fakejobId",
					Queued:   1,
					Complete: 0,
					Failed:   1,
					Running:  0,
					Starting: 0,
					Lost:     0,
				},
				{
					Name:     "Mandalorian",
					JobID:    "fakejobId",
					Queued:   1,
					Complete: 0,
					Failed:   0,
					Running:  1,
					Starting: 0,
					Lost:     0,
				},
			},
			Allocations: []*models.Alloc{},
			TaskGroupStatus: []*models.TaskGroupStatus{
				{
					ID:                "fakeID",
					Healthy:           1,
					Unhealthy:         0,
					Desired:           0,
					Placed:            0,
					ProgressDeadline:  0,
					Status:            "",
					StatusDescription: "",
				},
			},
		}

		//check that no error occured
		r.NoError(err)

		r.Equal(jId, "fakeID")

		//check that List() was called once
		r.Equal(fakeJobClient.InfoCallCount(), 1)

		//check that the query params where passed correctly
		r.Nil(queryOptions)

		//check all fields are set
		r.Equal(expectedJobStatus, jobStatus)
	})

	t.Run("When there is a problem with the client", func(t *testing.T) {
		fakeJobClient.InfoReturns(nil, nil, errors.New("aaah"))
		_, err := client.JobStatus("", nil)

		r.Error(err)
		r.Contains(err.Error(), "failed to retrieve job info")
	})
}
