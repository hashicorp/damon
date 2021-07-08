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

func TestDeployments(t *testing.T) {
	r := require.New(t)

	fakeClient := &nomadfakes.FakeDeploymentClient{}
	client := nomad.Nomad{DpClient: fakeClient}

	t.Run("When there are no issues", func(t *testing.T) {
		expectedDeps := []*models.Deployment{
			{
				ID:                "42",
				JobID:             "bumblebee",
				Namespace:         "transformers",
				Status:            "transformed",
				StatusDescription: "car",
			},
			{
				ID:                "23",
				JobID:             "optimus",
				Namespace:         "transformers",
				Status:            "transformed",
				StatusDescription: "truck",
			},
		}

		fakeClient.ListReturns([]*api.Deployment{
			{
				ID:                "42",
				JobID:             "bumblebee",
				Namespace:         "transformers",
				Status:            "transformed",
				StatusDescription: "car",
			},
			{
				ID:                "23",
				JobID:             "optimus",
				Namespace:         "transformers",
				Status:            "transformed",
				StatusDescription: "truck",
			},
		}, &api.QueryMeta{}, nil)

		dps, err := client.Deployments(&nomad.SearchOptions{
			Namespace: "transformers",
		})
		r.NoError(err)

		r.Equal(expectedDeps, dps)
	})

	t.Run("When there are no search options provided, it doesn't error", func(t *testing.T) {
		fakeClient.ListReturns([]*api.Deployment{}, &api.QueryMeta{}, nil)
		_, err := client.Deployments(nil)
		r.NoError(err)
	})

	t.Run("When there are issues with the client", func(t *testing.T) {
		fakeClient.ListReturns(nil, nil, errors.New("fatal"))
		_, err := client.Deployments(&nomad.SearchOptions{
			Namespace: "transformers",
		})
		r.Error(err)
		r.EqualError(err, "fatal")
	})
}
