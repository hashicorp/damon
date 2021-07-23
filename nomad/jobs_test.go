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

func TestJobs(t *testing.T) {
	r := require.New(t)

	fakeJobClient := &nomadfakes.FakeJobClient{}
	client := &nomad.Nomad{JobClient: fakeJobClient}

	t.Run("When there are no issues", func(t *testing.T) {
		now := time.Now().UnixNano()
		nowUnix := time.Unix(0, now)
		fakeJobClient.ListReturns([]*api.JobListStub{
			{
				ID:                "fake-id-1",
				Name:              "fake",
				Type:              "service",
				Status:            "running",
				StatusDescription: "this is awesome",
				SubmitTime:        now,
				JobSummary: &api.JobSummary{
					Namespace: "default",
				},
			},
			{
				ID:   "fake-id-2",
				Name: "fake-2",
				JobSummary: &api.JobSummary{
					Namespace: "default",
				},
			},
		}, nil, nil)

		so := nomad.SearchOptions{Namespace: "default"}

		jobs, err := client.Jobs(&so)

		queryOptions := fakeJobClient.ListArgsForCall(0)

		expectedJob := &models.Job{
			ID:                "fake-id-1",
			Name:              "fake",
			Namespace:         "default",
			Type:              "service",
			Status:            "running",
			StatusDescription: "this is awesome",
			SubmitTime:        nowUnix,
		}

		expectedQueryOptions := &api.QueryOptions{Namespace: "default"}

		//check that no error occured
		r.NoError(err)

		//check that List() was called once
		r.Equal(fakeJobClient.ListCallCount(), 1)

		//check that the query params where passed correctly
		r.Equal(queryOptions, expectedQueryOptions)

		//check that we have the right number of jobs
		r.Equal(len(jobs), 2)

		//check all fields are set
		r.Equal(jobs[0], expectedJob)
	})

	t.Run("When there is a problem with the client", func(t *testing.T) {
		fakeJobClient.ListReturns(nil, nil, errors.New("aaah"))
		_, err := client.Jobs(nil)

		r.Error(err)
		r.Contains(err.Error(), "failed to retrieve job list")
	})
}
