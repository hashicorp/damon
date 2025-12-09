// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad_test

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
)

func TestLogs(t *testing.T) {
	r := require.New(t)

	fakeFSClient := &nomadfakes.FakeAllocFSClient{}
	client := &nomad.Nomad{AllocFSClient: fakeFSClient}

	t.Run("It provides the correct params", func(t *testing.T) {
		allocID := "moon"
		taskName := "solar-system"
		logType := "stderr"
		cancelCh := make(<-chan struct{})

		client.Logs(allocID, taskName, logType, cancelCh)

		alloc,
			doFollow,
			actualTaskName,
			actualLogType,
			origin, offset,
			actualCancelCh,
			queryOptions := fakeFSClient.LogsArgsForCall(0)

		r.Equal(alloc, &api.Allocation{ID: "moon"})
		r.Equal(doFollow, true)
		r.Equal(actualTaskName, "solar-system")
		r.Equal(actualLogType, "stderr")
		r.Equal(origin, "end")
		r.Equal(offset, int64(20000))
		r.Equal(actualCancelCh, cancelCh)
		r.Equal(queryOptions, &api.QueryOptions{})
	})

	t.Run("It returns two channels", func(t *testing.T) {
		allocID := "moon"
		taskName := "solar-system"
		logType := "stderr"
		cancelCh := make(<-chan struct{})

		streamChan := make(<-chan *api.StreamFrame)
		errorChan := make(<-chan error)

		fakeFSClient.LogsReturns(streamChan, errorChan)
		actualStreamChan, actualErrorChan := client.Logs(allocID, taskName, logType, cancelCh)

		r.Equal(actualStreamChan, streamChan)
		r.Equal(actualErrorChan, errorChan)
	})
}
