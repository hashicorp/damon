// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
)

func TestExec(t *testing.T) {
	r := require.New(t)

	fakeAllocClient := &nomadfakes.FakeAllocationsClient{}
	client := &nomad.Nomad{AllocClient: fakeAllocClient}

	t.Run("It provides the correct params and returns the exit code", func(t *testing.T) {
		fakeAllocClient.ExecReturns(0, nil)

		ctx := context.Background()
		stdin := strings.NewReader("")
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		resizeCh := make(<-chan api.TerminalSize)
		command := []string{"/bin/sh"}

		exitCode, err := client.Exec(ctx, "moon", "solar-system", command, stdin, stdout, stderr, resizeCh)
		r.NoError(err)
		r.Equal(0, exitCode)

		actualCtx,
			alloc,
			taskName,
			tty,
			actualCommand,
			actualStdin,
			actualStdout,
			actualStderr,
			actualResizeCh,
			queryOptions := fakeAllocClient.ExecArgsForCall(0)

		r.Equal(ctx, actualCtx)
		r.Equal(&api.Allocation{ID: "moon"}, alloc)
		r.Equal("solar-system", taskName)
		r.True(tty)
		r.Equal(command, actualCommand)
		r.Equal(stdin, actualStdin)
		r.Equal(stdout, actualStdout)
		r.Equal(stderr, actualStderr)
		r.Equal(resizeCh, actualResizeCh)
		r.Equal(&api.QueryOptions{}, queryOptions)
	})

	t.Run("When the client fails", func(t *testing.T) {
		fakeAllocClient.ExecReturns(1, errors.New("argh"))

		exitCode, err := client.Exec(context.Background(), "moon", "solar-system",
			[]string{"/bin/sh"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}, nil)

		r.Error(err)
		r.EqualError(err, "argh")
		r.Equal(1, exitCode)
	})
}
