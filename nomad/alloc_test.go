package nomad_test

import (
	"errors"
	"testing"

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
				TaskStates: map[string]*api.TaskState{
					"task-1": {},
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
				TaskStates: map[string]*api.TaskState{
					"task-2": {},
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
				TaskNames:     []string{"task-1"},
			},
			{
				ID:            "id-two",
				TaskGroup:     "the-group",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
				TaskNames:     []string{"task-2"},
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
			},
			{
				ID:            "id-two",
				TaskGroup:     "the-group",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
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
			},
			{
				ID:            "id-two",
				TaskGroup:     "the-group",
				JobID:         "the-job",
				JobType:       "the-type",
				NodeID:        "node-id",
				NodeName:      "nodejs",
				DesiredStatus: "skate",
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
