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

func TestDeploymentTable_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		dt := component.NewDeploymentTable()

		dt.Table = fakeTable
		dt.Props.Namespace = "Namespace"
		dt.Props.Data = []*models.Deployment{
			{
				ID:                "ichi",
				JobID:             "saturn",
				Namespace:         "space",
				Status:            "running",
				StatusDescription: "fine",
			},
			{
				ID:                "ni",
				JobID:             "neptun",
				Namespace:         "outerspace",
				Status:            "failed",
				StatusDescription: "still fine",
			},
			{
				ID:                "san",
				JobID:             "jupiter",
				Namespace:         "outerspace",
				Status:            "pending",
				StatusDescription: "wait",
			},
			{
				ID:                "chi",
				JobID:             "pluto",
				Namespace:         "outerspace",
				Status:            "completed",
				StatusDescription: "...",
			},
		}

		dt.Props.SelectDeployment = func(id string) {}
		dt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		dt.Bind(slot)

		// It doesn't error
		err := dt.Render()
		r.NoError(err)

		// It renders the correct number of header rows
		renderHeaderCount := fakeTable.RenderHeaderCallCount()
		r.Equal(renderHeaderCount, 1)

		// It renders the correct header values
		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.TableHeaderDeployments, header)

		// It renders the correct number of rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 4)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		row2, index2, c2 := fakeTable.RenderRowArgsForCall(1)
		row3, index3, c3 := fakeTable.RenderRowArgsForCall(2)
		row4, index4, c4 := fakeTable.RenderRowArgsForCall(3)

		expectedRow1 := []string{"ichi", "saturn", "space", "running", "fine"}
		expectedRow2 := []string{"ni", "neptun", "outerspace", "failed", "still fine"}
		expectedRow3 := []string{"san", "jupiter", "outerspace", "pending", "wait"}
		expectedRow4 := []string{"chi", "pluto", "outerspace", "completed", "..."}

		// It render the correct data for the rows
		r.Equal(expectedRow1, row1)
		r.Equal(expectedRow2, row2)
		r.Equal(expectedRow3, row3)
		r.Equal(expectedRow4, row4)

		// It renders the data at the correct index
		r.Equal(index1, 1)
		r.Equal(index2, 2)
		r.Equal(index3, 3)
		r.Equal(index4, 4)

		// It renders the rows in the correct color
		r.Equal(c1, styles.TcellColorHighlighPrimary)
		r.Equal(c2, tcell.ColorRed)
		r.Equal(c3, tcell.ColorYellow)
		r.Equal(c4, tcell.ColorWhite)
	})

	t.Run("When render called again", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		dt := component.NewDeploymentTable()

		dt.Table = fakeTable
		dt.Props.Namespace = "Namespace"
		dt.Props.Data = []*models.Deployment{
			{
				ID:                "ichi",
				JobID:             "saturn",
				Namespace:         "space",
				Status:            "running",
				StatusDescription: "fine",
			},
		}

		dt.Props.SelectDeployment = func(id string) {}
		dt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		dt.Bind(slot)

		// It doesn't error
		err := dt.Render()
		r.NoError(err)

		err = dt.Render()
		r.NoError(err)

		// It clears the table on each call
		r.Equal(fakeTable.ClearCallCount(), 2)
	})

	t.Run("When there is no data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		dt := component.NewDeploymentTable()

		dt.Table = fakeTable
		dt.Props.Namespace = "Namespace"
		dt.Props.Data = []*models.Deployment{}

		dt.Props.SelectDeployment = func(id string) {}

		var handleNoResourcesCalled bool
		dt.Props.HandleNoResources = func(format string, args ...interface{}) {
			handleNoResourcesCalled = true

			r.Equal("%sno deployments available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}

		slot := tview.NewFlex()
		dt.Bind(slot)

		// It doesn't error
		err := dt.Render()
		r.NoError(err)

		// It handled the case that there are no resources
		r.True(handleNoResourcesCalled)

		// It didn't returned after handling no resources
		r.Equal(fakeTable.RenderHeaderCallCount(), 0)
		r.Equal(fakeTable.RenderRowCallCount(), 0)
	})
}

func TestDeploymentTable_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When SelectDeployment is not set", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		dt := component.NewDeploymentTable()

		dt.Table = fakeTable

		dt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		// It errors
		err := dt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When HandleNoResources is not set", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		dt := component.NewDeploymentTable()

		dt.Table = fakeTable

		dt.Props.SelectDeployment = func(id string) {}

		// It errors
		err := dt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When the component isn't bound", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		dt := component.NewDeploymentTable()

		dt.Table = fakeTable

		dt.Props.SelectDeployment = func(id string) {}
		dt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		// It errors
		err := dt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})
}
