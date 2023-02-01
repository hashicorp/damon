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
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/styles"
)

func TestTaskGroupTable_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		tt := component.NewTaskGroupTable()

		tt.Table = fakeTable
		tt.Props.Data = []*models.TaskGroup{
			{
				Name:  "ichi",
				JobID: "earth",
			},
			{
				Name:  "ni",
				JobID: "neptun",
			},
		}

		tt.Props.SelectTaskGroup = func(id string) {}
		tt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		tt.Bind(slot)

		// It doesn't error
		err := tt.Render()
		r.NoError(err)

		// It renders the correct number of header rows
		renderHeaderCount := fakeTable.RenderHeaderCallCount()
		r.Equal(renderHeaderCount, 1)

		// It renders the correct header values
		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.TableHeaderTaskGroups, header)

		// It renders the correct number of rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 2)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		row2, index2, c2 := fakeTable.RenderRowArgsForCall(1)

		expectedRow1 := []string{"ichi", "earth", "0", "0", "0", "0", "0", "0"}
		expectedRow2 := []string{"ni", "neptun", "0", "0", "0", "0", "0", "0"}

		// It render the correct data for the rows
		r.Equal(expectedRow1, row1)
		r.Equal(expectedRow2, row2)

		// It renders the data at the correct index
		r.Equal(index1, 1)
		r.Equal(index2, 2)

		// It renders the rows in the correct color
		r.Equal(c1, tcell.ColorWhite)
		r.Equal(c2, tcell.ColorWhite)
	})

	t.Run("When render called again", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		tt := component.NewTaskGroupTable()

		tt.Table = fakeTable
		tt.Props.Data = []*models.TaskGroup{
			{
				Name:  "ichi",
				JobID: "earth",
			},
		}

		tt.Props.SelectTaskGroup = func(id string) {}
		tt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		tt.Bind(slot)

		// It doesn't error
		err := tt.Render()
		r.NoError(err)

		err = tt.Render()
		r.NoError(err)

		// It clears the table on each call
		r.Equal(fakeTable.ClearCallCount(), 2)
	})

	t.Run("When there is no data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		tt := component.NewTaskGroupTable()

		tt.Table = fakeTable
		tt.Props.Data = []*models.TaskGroup{}

		tt.Props.SelectTaskGroup = func(id string) {}
		var handleNoResourcesCalled bool
		tt.Props.HandleNoResources = func(format string, args ...interface{}) {
			handleNoResourcesCalled = true

			r.Equal("%sno TaskGroups available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}

		slot := tview.NewFlex()
		tt.Bind(slot)

		// It doesn't error
		err := tt.Render()
		r.NoError(err)

		// It handled the case that there are no resources
		r.True(handleNoResourcesCalled)

		// It didn't returned after handling no resources
		r.Equal(fakeTable.RenderHeaderCallCount(), 0)
		r.Equal(fakeTable.RenderRowCallCount(), 0)
	})
}

func TestTaskGroupTable_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When SelectTaskGroup is not set", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		tt := component.NewTaskGroupTable()

		tt.Table = fakeTable
		tt.Props.Data = []*models.TaskGroup{}

		tt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		tt.Bind(slot)

		// It doesn't error
		err := tt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When HandleNoResources is not set", func(t *testing.T) {
		tt := component.NewTaskGroupTable()
		tt.Props.SelectTaskGroup = func(id string) {}

		slot := tview.NewFlex()
		tt.Bind(slot)

		// It doesn't error
		err := tt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When the component isn't bound", func(t *testing.T) {
		tt := component.NewTaskGroupTable()

		tt.Props.SelectTaskGroup = func(id string) {}
		tt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		// It doesn't error
		err := tt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})
}
