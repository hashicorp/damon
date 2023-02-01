// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
)

func TestSearch_Happy(t *testing.T) {
	r := require.New(t)

	input := &componentfakes.FakeInputField{}
	search := component.NewSearchField("test")
	search.InputField = input

	var changedCalled bool
	search.Props.ChangedFunc = func(text string) {
		changedCalled = true
	}

	var doneCalled bool
	search.Props.DoneFunc = func(key tcell.Key) {
		doneCalled = true
	}
	search.Bind(tview.NewFlex())

	err := search.Render()
	r.NoError(err)

	actualDoneFunc := input.SetDoneFuncArgsForCall(0)
	actualChangedFunc := input.SetChangedFuncArgsForCall(0)

	actualChangedFunc("")
	actualDoneFunc(tcell.KeyACK)

	r.True(changedCalled)
	r.True(doneCalled)
}

func TestSearch_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the component isn't bound", func(t *testing.T) {
		input := &componentfakes.FakeInputField{}
		search := component.NewSearchField("test")
		search.InputField = input
		search.Props.ChangedFunc = func(text string) {}
		search.Props.DoneFunc = func(key tcell.Key) {}

		err := search.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})

	t.Run("When DoneFunc is not set", func(t *testing.T) {
		input := &componentfakes.FakeInputField{}
		search := component.NewSearchField("test")
		search.InputField = input
		search.Props.ChangedFunc = func(text string) {}
		search.Bind(tview.NewFlex())

		err := search.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When ChangedFunc is not set", func(t *testing.T) {
		input := &componentfakes.FakeInputField{}
		search := component.NewSearchField("test")
		search.InputField = input
		search.Props.DoneFunc = func(key tcell.Key) {}
		search.Bind(tview.NewFlex())

		err := search.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})
}
