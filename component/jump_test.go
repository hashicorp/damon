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

func TestJump_Happy(t *testing.T) {
	r := require.New(t)

	input := &componentfakes.FakeInputField{}
	jump := component.NewJumpToJob()
	jump.InputField = input

	var doneCalled bool
	jump.Props.DoneFunc = func(key tcell.Key) {
		doneCalled = true
	}

	jump.Bind(tview.NewFlex())

	err := jump.Render()
	r.NoError(err)

	actualDoneFunc := input.SetDoneFuncArgsForCall(0)

	actualDoneFunc(tcell.KeyACK)

	r.True(doneCalled)
}

func TestJump_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the component isn't bound", func(t *testing.T) {
		jump := component.NewJumpToJob()

		jump.Props.DoneFunc = func(key tcell.Key) {}

		err := jump.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})

	t.Run("When DoneFunc is not set", func(t *testing.T) {
		jump := component.NewJumpToJob()

		jump.Bind(tview.NewFlex())

		err := jump.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

}
