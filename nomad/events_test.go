// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
)

func TestStream(t *testing.T) {
	r := require.New(t)

	client := &nomadfakes.FakeEventsClient{}
	nmd := &nomad.Nomad{EventsClient: client}

	topics := map[api.Topic][]string{
		"Job": {"*"},
	}

	t.Run("It provides the correct params", func(t *testing.T) {
		topics := map[api.Topic][]string{
			"Job": {"*"},
		}

		_, err := nmd.Stream(topics, 0)
		r.NoError(err)

	})

	t.Run("It returns a channel and an error", func(t *testing.T) {
		streamChan := make(<-chan *api.Events)
		err := errors.New("haha")

		client.StreamReturns(streamChan, err)
		actualStreamChan, actualError := nmd.Stream(topics, 0)

		r.Equal(actualStreamChan, streamChan)
		r.Equal(actualError, err)
	})
}
