// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
	"github.com/hcjulz/damon/state"
)

func newFailureView() (*View, *componentfakes.FakeModal) {
	modal := &componentfakes.FakeModal{}

	failure := component.NewInfo()
	failure.Modal = modal
	failure.Props.Done = func(buttonIndex int, buttonLabel string) {}
	failure.Bind(tview.NewPages())

	view := New(&Components{Failure: failure}, nil, nil, &state.State{})

	return view, modal
}

func TestErr(t *testing.T) {
	t.Run("When the message contains printf verbs", func(t *testing.T) {
		r := require.New(t)

		view, modal := newFailureView()

		// handleError formats its arguments, so a message carrying verbs must
		// reach the modal unchanged.
		view.err(errors.New("connection refused"), "cannot reach 100% of nodes: %s")

		r.Equal(1, modal.SetTextCallCount())
		r.Equal("cannot reach 100% of nodes: %s", modal.SetTextArgsForCall(0))
	})

	t.Run("When there is no error", func(t *testing.T) {
		r := require.New(t)

		view, modal := newFailureView()

		view.err(nil, "no failure happened")

		r.Equal(0, modal.SetTextCallCount())
	})
}

func TestHandleError(t *testing.T) {
	r := require.New(t)

	view, modal := newFailureView()

	view.handleError("%s is %s", "nomad", "unreachable")

	r.Equal(1, modal.SetTextCallCount())
	r.Equal("nomad is unreachable", modal.SetTextArgsForCall(0))
}
