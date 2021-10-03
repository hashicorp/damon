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
	client := &nomad.Nomad{JobClient: fakeClient}

	now := time.Now()

	t.Run("When there are no issues", func(t *testing.T) {
		fakeClient.AllocationsReturns([]*api.AllocationListStub{
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

		qo := &nomad.SearchOptions{
			Namespace: "nodejs",
		}

		allocs, err := client.JobAllocs("the-job", qo)
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
