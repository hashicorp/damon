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
					Summary: map[string]api.TaskGroupSummary{
						"task1": {Running: 0},
						"task2": {Running: 1},
					},
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
			StatusSummary:     models.Summary{Total: 2, Running: 1},
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

func TestStartJob(t *testing.T) {
	r := require.New(t)

	fakeJobClient := &nomadfakes.FakeJobClient{}
	client := &nomad.Nomad{JobClient: fakeJobClient}

	t.Run("When everything is fine", func(t *testing.T) {
		fakeJobClient.RegisterReturns(&api.JobRegisterResponse{}, &api.WriteMeta{}, nil)

		id := "test"
		job := api.Job{ID: &id}
		err := client.StartJob(&job)
		r.NoError(err)

		actualJob, writeOpts := fakeJobClient.RegisterArgsForCall(0)

		r.Equal(actualJob, &job)
		r.Nil(writeOpts)
	})

	t.Run("When the client is failing", func(t *testing.T) {
		fakeJobClient.RegisterReturns(nil, nil, errors.New("argh"))

		id := "test"
		job := api.Job{ID: &id}
		err := client.StartJob(&job)

		r.Error(err)
		r.EqualError(err, "argh")
	})
}

func TestStopJob(t *testing.T) {
	r := require.New(t)

	fakeJobClient := &nomadfakes.FakeJobClient{}
	client := &nomad.Nomad{JobClient: fakeJobClient}

	t.Run("When everything is fine", func(t *testing.T) {
		fakeJobClient.DeregisterReturns("test", &api.WriteMeta{}, nil)

		err := client.StopJob("test")
		r.NoError(err)

		actualJobID, purge, writeOpts := fakeJobClient.DeregisterArgsForCall(0)

		r.Equal(actualJobID, "test")
		r.False(purge)
		r.Nil(writeOpts)
	})

	t.Run("When the client is failing", func(t *testing.T) {
		fakeJobClient.DeregisterReturns("", nil, errors.New("argh"))

		err := client.StopJob("test")
		r.Error(err)
		r.EqualError(err, "argh")
	})
}

func TestGetJob(t *testing.T) {
	r := require.New(t)

	fakeJobClient := &nomadfakes.FakeJobClient{}
	client := &nomad.Nomad{JobClient: fakeJobClient}

	t.Run("When everything is fine", func(t *testing.T) {
		id := "test"
		fakeJobClient.InfoReturns(&api.Job{ID: &id}, nil, nil)
		job, err := client.GetJob("test")
		r.NoError(err)

		actualJobID, queryOptions := fakeJobClient.InfoArgsForCall(0)

		r.Equal(actualJobID, "test")
		r.Equal(*job.ID, "test")
		r.Nil(queryOptions)
	})

	t.Run("When the client is failing", func(t *testing.T) {
		fakeJobClient.InfoReturns(nil, nil, errors.New("argh"))

		_, err := client.GetJob("test")
		r.Error(err)
		r.EqualError(err, "argh")
	})
}
