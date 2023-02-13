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
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/state"
)

func TestSelections_Happy(t *testing.T) {
	r := require.New(t)

	state := state.New()
	state.Namespaces = []*models.Namespace{
		{
			Name:        "test",
			Description: "test-space",
		},
		{
			Name:        "space",
			Description: "ship",
		},
	}
	dropdown := &componentfakes.FakeDropDown{}

	selections := component.NewSelections(state)
	selections.Namespace = dropdown

	selections.Bind(tview.NewFlex())

	dropdown.PrimitiveReturns(tview.NewDropDown())

	err := selections.Render()
	r.NoError(err)

	ns, _ := dropdown.SetOptionsArgsForCall(0)
	optIndex := dropdown.SetCurrentOptionArgsForCall(0)
	actualRerender := dropdown.SetSelectedFuncArgsForCall(0)

	r.Equal(ns, []string{"test", "space"})
	r.Equal(optIndex, 1)

	actualRerender("text", 0)
}

func TestSelections_Sad(t *testing.T) {
	t.Run("When the component isn't bound", func(t *testing.T) {
		r := require.New(t)

		state := state.New()
		state.Namespaces = []*models.Namespace{}
		dropdown := &componentfakes.FakeDropDown{}

		selections := component.NewSelections(state)
		selections.Namespace = dropdown

		dropdown.PrimitiveReturns(tview.NewDropDown())

		err := selections.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})
}
