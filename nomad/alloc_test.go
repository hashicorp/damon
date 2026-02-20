// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad_test

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
)

func TestJobAllocs(t *testing.T) {
	r := require.New(t)

	fakeClient := &nomadfakes.FakeJobClient{}
	fakeAllocClient := &nomadfakes.FakeAllocationsClient{}
	client := &nomad.Nomad{
		JobClient:   fakeClient,
		AllocClient: fakeAllocClient,
	}

	now := time.Now()

	t.Run("When there are no issues", func(t *testing.T) {
		fakeClient.AllocationsReturns([]*api.AllocationListStub{
			{
				ID:            "id-one",
				TaskGroup:     "the-group",
				Namespace:     "namespace",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				JobVersion:    100,
				ClientStatus:  "ClientStatus",
				CreateTime:    100,
				ModifyTime:    100,
				TaskStates: map[string]*api.TaskState{
					"task-1": {
						State: "running",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
			{
				ID:            "id-two",
				TaskGroup:     "the-group",
				Namespace:     "namespace",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				JobVersion:    200,
				ClientStatus:  "ClientStatus",
				CreateTime:    200,
				ModifyTime:    200,
				TaskStates: map[string]*api.TaskState{
					"task-2": {
						State: "running",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
			{
				ID:            "id-three",
				TaskGroup:     "the-group",
				Namespace:     "namespace",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				JobVersion:    200,
				ClientStatus:  "ClientStatus",
				CreateTime:    200,
				ModifyTime:    200,
				TaskStates:    map[string]*api.TaskState{},
			},
		}, &api.QueryMeta{}, nil)

		cpu, memory := 100, 10
		tgName := "the-group"
		fakeAllocClient.InfoReturnsOnCall(0, &api.Allocation{
			TaskGroup: "the-group",
			Job: &api.Job{
				TaskGroups: []*api.TaskGroup{
					{
						Name: &tgName,
						Tasks: []*api.Task{
							{
								Name:   "task-1",
								Driver: "docker",
								Env: map[string]string{
									"env-key": "env-value",
								},
								Config: map[string]interface{}{
									"image": "the-image-i-run",
								},
								Resources: &api.Resources{
									CPU:      &cpu,
									MemoryMB: &memory,
								},
							},
						},
					},
				},
			},
		}, &api.QueryMeta{}, nil)

		fakeAllocClient.InfoReturnsOnCall(1, &api.Allocation{
			TaskGroup: "the-group",
			Job: &api.Job{
				TaskGroups: []*api.TaskGroup{
					{
						Name: &tgName,
						Tasks: []*api.Task{
							{
								Name:   "task-2",
								Driver: "docker",
								Env: map[string]string{
									"env-key": "env-value",
								},
								Config: map[string]interface{}{
									"image": "the-image-i-run",
								},
								Resources: &api.Resources{
									CPU:      &cpu,
									MemoryMB: &memory,
								},
							},
						},
					},
				},
			},
		}, &api.QueryMeta{}, nil)

		fakeAllocClient.InfoReturnsOnCall(2, &api.Allocation{
			TaskGroup: "the-group",
			Job: &api.Job{
				TaskGroups: []*api.TaskGroup{
					{
						Name: &tgName,
						Tasks: []*api.Task{
							{
								Name:   "task-3",
								Driver: "docker",
								Env: map[string]string{
									"env-key": "env-value",
								},
								Config: map[string]interface{}{
									"image": "the-image-i-run",
								},
								Resources: &api.Resources{
									CPU:      &cpu,
									MemoryMB: &memory,
								},
							},
						},
					},
				},
			},
		}, &api.QueryMeta{}, nil)

		qo := &nomad.SearchOptions{
			Namespace: "namespace",
		}

		allocs, err := client.JobAllocs("the-job", qo)
		r.NoError(err)

		expectedAllocs := []*models.Alloc{
			{
				ID:            "id-one",
				TaskGroup:     "the-group",
				Namespace:     "namespace",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				Version:       100,
				Status:        "ClientStatus",
				Created:       time.Unix(0, 100),
				Modified:      time.Unix(0, 100),
				TaskNames:     []string{"task-1"},
				Tasks: []models.AllocTask{
					{
						Name: "task-1",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
				TaskList: []*models.Task{
					{
						Name:   "task-1",
						Driver: "docker",
						State:  "running",
						Env: map[string]string{
							"env-key": "env-value",
						},
						Config: map[string]interface{}{
							"image": "the-image-i-run",
						},
						CPU:      cpu,
						MemoryMB: memory,
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
			{
				ID:            "id-two",
				TaskGroup:     "the-group",
				Namespace:     "namespace",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				Version:       200,
				Status:        "ClientStatus",
				Created:       time.Unix(0, 200),
				Modified:      time.Unix(0, 200),
				TaskNames:     []string{"task-2"},
				Tasks: []models.AllocTask{
					{
						Name: "task-2",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
				TaskList: []*models.Task{
					{
						Name:   "task-2",
						Driver: "docker",
						State:  "running",
						Env: map[string]string{
							"env-key": "env-value",
						},
						Config: map[string]interface{}{
							"image": "the-image-i-run",
						},
						CPU:      cpu,
						MemoryMB: memory,
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
			{
				ID:            "id-three",
				TaskGroup:     "the-group",
				Namespace:     "namespace",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				Version:       200,
				Status:        "ClientStatus",
				Created:       time.Unix(0, 200),
				Modified:      time.Unix(0, 200),
				TaskNames:     []string{"task-3"},
				Tasks: []models.AllocTask{
					{
						Name: "task-3",
					},
				},
				TaskList: []*models.Task{
					{
						Name:   "task-3",
						Driver: "docker",
						Env: map[string]string{
							"env-key": "env-value",
						},
						Config: map[string]interface{}{
							"image": "the-image-i-run",
						},
						CPU:      cpu,
						MemoryMB: memory,
					},
				},
			},
		}

		jobID, all, apiQo := fakeClient.AllocationsArgsForCall(0)

		r.Equal(fakeClient.AllocationsCallCount(), 1)

		r.Equal(expectedAllocs, allocs)
		r.Equal(jobID, "the-job")
		r.Equal(all, false)
		r.Equal(apiQo, &api.QueryOptions{
			Namespace: qo.Namespace,
		})
	})

	t.Run("When there are no SearchOptions provided", func(t *testing.T) {
		fakeClient.AllocationsReturns([]*api.AllocationListStub{}, &api.QueryMeta{}, nil)
		_, err := client.JobAllocs("id", nil)
		r.NoError(err)
	})

	t.Run("When there is an issue with the client", func(t *testing.T) {
		fakeClient.AllocationsReturns(nil, nil, errors.New("argh"))
		_, err := client.JobAllocs("id", nil)

		r.Error(err)
		r.EqualError(err, "argh")
	})
}

func TestAllocations(t *testing.T) {
	r := require.New(t)

	fakeClient := &nomadfakes.FakeAllocationsClient{}
	client := &nomad.Nomad{AllocClient: fakeClient}

	now := time.Now()

	t.Run("When there are no issues", func(t *testing.T) {
		fakeClient.ListReturns([]*api.AllocationListStub{
			{
				ID:            "id-one",
				TaskGroup:     "the-group",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				JobVersion:    100,
				ClientStatus:  "ClientStatus",
				CreateTime:    100,
				ModifyTime:    100,
				TaskStates: map[string]*api.TaskState{
					"task-1": {
						State: "running",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
			{
				ID:            "id-two",
				TaskGroup:     "the-group",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				JobVersion:    200,
				ClientStatus:  "ClientStatus",
				CreateTime:    200,
				ModifyTime:    200,
				TaskStates: map[string]*api.TaskState{
					"task-2": {
						State: "running",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
		}, &api.QueryMeta{}, nil)

		cpu, memory := 100, 10
		tgName := "the-group"
		fakeClient.InfoReturnsOnCall(0, &api.Allocation{
			TaskGroup: "the-group",
			Job: &api.Job{
				TaskGroups: []*api.TaskGroup{
					{
						Name: &tgName,
						Tasks: []*api.Task{
							{
								Name:   "task-1",
								Driver: "docker",
								Env: map[string]string{
									"env-key": "env-value",
								},
								Config: map[string]interface{}{
									"image": "the-image-i-run",
								},
								Resources: &api.Resources{
									CPU:      &cpu,
									MemoryMB: &memory,
								},
							},
						},
					},
				},
			},
		}, &api.QueryMeta{}, nil)

		fakeClient.InfoReturnsOnCall(1, &api.Allocation{
			TaskGroup: "the-group",
			Job: &api.Job{
				TaskGroups: []*api.TaskGroup{
					{
						Name: &tgName,
						Tasks: []*api.Task{
							{
								Name:   "task-2",
								Driver: "docker",
								Env: map[string]string{
									"env-key": "env-value",
								},
								Config: map[string]interface{}{
									"image": "the-image-i-run",
								},
								Resources: &api.Resources{
									CPU:      &cpu,
									MemoryMB: &memory,
								},
							},
						},
					},
				},
			},
		}, &api.QueryMeta{}, nil)

		qo := &nomad.SearchOptions{
			Namespace: "*",
		}

		allocs, err := client.Allocations(qo)
		r.NoError(err)

		expectedAllocs := []*models.Alloc{
			{
				ID:            "id-one",
				TaskGroup:     "the-group",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				Version:       100,
				Status:        "ClientStatus",
				Created:       time.Unix(0, 100),
				Modified:      time.Unix(0, 100),
				TaskNames:     []string{"task-1"},
				Tasks: []models.AllocTask{
					{
						Name: "task-1",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
				TaskList: []*models.Task{
					{
						Name:   "task-1",
						Driver: "docker",
						State:  "running",
						Env: map[string]string{
							"env-key": "env-value",
						},
						Config: map[string]interface{}{
							"image": "the-image-i-run",
						},
						CPU:      cpu,
						MemoryMB: memory,
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
			{
				ID:            "id-two",
				TaskGroup:     "the-group",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				Version:       200,
				Status:        "ClientStatus",
				Created:       time.Unix(0, 200),
				Modified:      time.Unix(0, 200),
				TaskNames:     []string{"task-2"},
				Tasks: []models.AllocTask{
					{
						Name: "task-2",
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
				TaskList: []*models.Task{
					{
						Name:   "task-2",
						Driver: "docker",
						State:  "running",
						Env: map[string]string{
							"env-key": "env-value",
						},
						Config: map[string]interface{}{
							"image": "the-image-i-run",
						},
						CPU:      cpu,
						MemoryMB: memory,
						Events: []*api.TaskEvent{
							{
								Time:           now.UnixNano(),
								Type:           "type",
								DisplayMessage: "msg",
							},
						},
					},
				},
			},
		}

		apiQo := fakeClient.ListArgsForCall(0)

		r.Equal(expectedAllocs, allocs)
		r.Equal(apiQo, &api.QueryOptions{
			Namespace: qo.Namespace,
		})
	})

	t.Run("When there are no SearchOptions provided", func(t *testing.T) {
		fakeClient.ListReturns([]*api.AllocationListStub{}, &api.QueryMeta{}, nil)
		_, err := client.Allocations(nil)
		r.NoError(err)
	})

	t.Run("When there is an issue with the client", func(t *testing.T) {
		fakeClient.ListReturns(nil, nil, errors.New("argh"))
		_, err := client.Allocations(nil)
		r.Error(err)
		r.EqualError(err, "argh")
	})
}
