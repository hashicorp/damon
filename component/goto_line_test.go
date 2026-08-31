// Copyright IBM Corp. 2021, 2023
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

func TestGotoLine_Happy(t *testing.T) {
	r := require.New(t)

	input := &componentfakes.FakeInputField{}
	gotoLine := component.NewGotoLine()
	gotoLine.InputField = input

	var doneCalled bool
	gotoLine.Props.DoneFunc = func(key tcell.Key) {
		doneCalled = true
	}

	gotoLine.Bind(tview.NewFlex())

	err := gotoLine.Render()
	r.NoError(err)

	actualDoneFunc := input.SetDoneFuncArgsForCall(0)

	actualDoneFunc(tcell.KeyEnter)

	r.True(doneCalled)
}

func TestGotoLine_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the component isn't bound", func(t *testing.T) {
		gotoLine := component.NewGotoLine()

		gotoLine.Props.DoneFunc = func(key tcell.Key) {}

		err := gotoLine.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})

	t.Run("When DoneFunc is not set", func(t *testing.T) {
		gotoLine := component.NewGotoLine()

		gotoLine.Bind(tview.NewFlex())

		err := gotoLine.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})
}
