// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
)

func TestError_Happy(t *testing.T) {
	r := require.New(t)

	e := component.NewError()

	modal := &componentfakes.FakeModal{}
	e.Modal = modal

	var doneCalled bool
	e.Props.Done = func(buttonIndex int, buttonLabel string) {
		doneCalled = true
	}

	pages := tview.NewPages()
	e.Bind(pages)

	err := e.Render("error")
	r.NoError(err)

	actualDone := modal.SetDoneFuncArgsForCall(0)
	text := modal.SetTextArgsForCall(0)

	actualDone(0, "buttonName")

	r.True(doneCalled)
	r.Equal(text, "error")
}

func TestError_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the component isn't bound", func(t *testing.T) {
		e := component.NewError()

		modal := &componentfakes.FakeModal{}
		e.Modal = modal

		e.Props.Done = func(buttonIndex int, buttonLabel string) {}

		err := e.Render("error")
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})

	t.Run("When DoneFunc is not set", func(t *testing.T) {
		e := component.NewError()

		modal := &componentfakes.FakeModal{}
		e.Modal = modal

		pages := tview.NewPages()
		e.Bind(pages)

		err := e.Render("error")
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})
}
