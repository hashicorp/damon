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

func TestTaskGroups_Happy(t *testing.T) {
	r := require.New(t)

	fakeJobClient := &nomadfakes.FakeJobClient{}
	client := &nomad.Nomad{JobClient: fakeJobClient}

	jobSummary := &api.JobSummary{
		JobID: "StarWars",
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
	}

	fakeJobClient.SummaryReturns(jobSummary, nil, nil)

	t.Run("It doesn't fail", func(t *testing.T) {
		_, err := client.TaskGroups("ID", nil)
		r.NoError(err)
	})

	t.Run("It returns the correct number of tasks", func(t *testing.T) {
		tg, _ := client.TaskGroups("ID", nil)
		r.Len(tg, 2)
	})

	t.Run("It returns the expected task groups", func(t *testing.T) {
		actualTaskGroups, _ := client.TaskGroups("ID", nil)

		expectedTaskGroups := []*models.TaskGroup{
			{
				Name:   "Grogu",
				JobID:  "StarWars",
				Failed: 1,
				Queued: 1,
			},
			{
				Name:    "Mandalorian",
				JobID:   "StarWars",
				Running: 1,
				Queued:  1,
			},
		}

		r.Equal(expectedTaskGroups, actualTaskGroups)
	})
}

func TestTaskGroups_Sad(t *testing.T) {
	r := require.New(t)

	fakeJobClient := &nomadfakes.FakeJobClient{}
	nomad := &nomad.Nomad{JobClient: fakeJobClient}
	fakeJobClient.SummaryReturns(nil, nil, errors.New("aaah"))

	t.Run("It fails when the client returns an error", func(t *testing.T) {
		_, err := nomad.TaskGroups("ID", nil)
		r.Error(err)
	})
}
