// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/styles"
)

func TestJobTable_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		jt := component.NewJobsTable()

		now := time.Now()

		jt.Table = fakeTable
		jt.Props.Namespace = "Namespace"
		jt.Props.Data = []*models.Job{
			{
				ID:                "ichi",
				Name:              "saturn",
				Type:              "service",
				Namespace:         "space",
				Status:            "running",
				StatusDescription: "fine",
				ReadyStatus:       models.ReadyStatus{Desired: 2, Running: 2, Healthy: 2},
				SubmitTime:        now,
			},
			{
				ID:                "ni",
				Name:              "jupiter",
				Type:              "service",
				Namespace:         "space",
				Status:            "pending",
				StatusDescription: "fine",
				ReadyStatus:       models.ReadyStatus{Desired: 2, Running: 2, Healthy: 2},
				SubmitTime:        now,
			},
			{
				ID:                "san",
				Name:              "neptun",
				Type:              "service",
				Namespace:         "space",
				Status:            "dead",
				StatusDescription: "fine",
				ReadyStatus:       models.ReadyStatus{Desired: 1, Running: 1, Healthy: 1},
				SubmitTime:        now,
			},
			{
				ID:                "chi",
				Name:              "mars",
				Type:              "batch",
				Namespace:         "space",
				Status:            "running",
				ReadyStatus:       models.ReadyStatus{Desired: 1, Running: 1, Healthy: 0},
				StatusDescription: "fine",
				SubmitTime:        now,
			},
			{
				ID:                "yo",
				Name:              "venus",
				Type:              "service",
				Namespace:         "space",
				Status:            "running",
				ReadyStatus:       models.ReadyStatus{Desired: 1, Running: 1, Healthy: 0, Unhealthy: 1},
				StatusDescription: "fine",
				SubmitTime:        now,
			},
		}

		jt.Props.SelectJob = func(id string) {}
		jt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		jt.Bind(slot)

		// It doesn't error
		err := jt.Render()
		r.NoError(err)

		// It renders the correct number of header rows
		renderHeaderCount := fakeTable.RenderHeaderCallCount()
		r.Equal(renderHeaderCount, 1)

		// It renders the correct header values
		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.TableHeaderJobs, header)

		// It renders the correct number of rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 5)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		row2, index2, c2 := fakeTable.RenderRowArgsForCall(1)
		row3, index3, c3 := fakeTable.RenderRowArgsForCall(2)
		row4, index4, c4 := fakeTable.RenderRowArgsForCall(3)
		row5, index5, c5 := fakeTable.RenderRowArgsForCall(4)

		expectedRow1 := []string{"1", "ichi", "saturn", "service", "space", "running", "2/2 ✅", now.Format(time.RFC3339), "0s"}
		expectedRow2 := []string{"2", "ni", "jupiter", "service", "space", "pending", "---", now.Format(time.RFC3339), "0s"}
		expectedRow3 := []string{"3", "san", "neptun", "service", "space", "dead", "---", now.Format(time.RFC3339), "0s"}
		expectedRow4 := []string{"4", "chi", "mars", "batch", "space", "running", "0/1 ⚠️", now.Format(time.RFC3339), "0s"}
		expectedRow5 := []string{"5", "yo", "venus", "service", "space", "running", "0/1 ❌", now.Format(time.RFC3339), "0s"}

		// It render the correct data for the rows
		r.Equal(expectedRow1, row1)
		r.Equal(expectedRow2, row2)
		r.Equal(expectedRow3, row3)
		r.Equal(expectedRow4, row4)
		r.Equal(expectedRow5, row5)

		// It renders the data at the correct index
		r.Equal(index1, 1)
		r.Equal(index2, 2)
		r.Equal(index3, 3)
		r.Equal(index4, 4)
		r.Equal(index5, 5)

		// It renders the rows in the correct color
		r.Equal(c1, tcell.ColorWhite)
		r.Equal(c2, tcell.ColorDarkGrey)
		r.Equal(c3, tcell.ColorDarkGrey)
		r.Equal(c4, tcell.ColorWhite)
		r.Equal(c5, tcell.ColorRed)
	})

	t.Run("When render called again", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		jt := component.NewJobsTable()

		now := time.Now()

		jt.Table = fakeTable
		jt.Props.Namespace = "Namespace"
		jt.Props.Data = []*models.Job{
			{
				ID:                "ichi",
				Name:              "saturn",
				Type:              "service",
				Namespace:         "space",
				Status:            "running",
				StatusDescription: "fine",
				StatusSummary:     models.Summary{Total: 1, Running: 1},
				SubmitTime:        now,
			},
		}

		jt.Props.SelectJob = func(id string) {}
		jt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		jt.Bind(slot)

		// It doesn't error
		err := jt.Render()
		r.NoError(err)

		err = jt.Render()
		r.NoError(err)

		// It clears the table on each call
		r.Equal(fakeTable.ClearCallCount(), 2)
	})

	t.Run("When there is no data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		jt := component.NewJobsTable()

		jt.Table = fakeTable
		jt.Props.Namespace = "Namespace"
		jt.Props.Data = []*models.Job{}

		jt.Props.SelectJob = func(id string) {}

		var handleNoResourcesCalled bool
		jt.Props.HandleNoResources = func(format string, args ...interface{}) {
			handleNoResourcesCalled = true

			r.Equal("%sno jobs available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}

		slot := tview.NewFlex()
		jt.Bind(slot)

		// It doesn't error
		err := jt.Render()
		r.NoError(err)

		// It handled the case that there are no resources
		r.True(handleNoResourcesCalled)

		// It didn't returned after handling no resources
		r.Equal(fakeTable.RenderHeaderCallCount(), 0)
		r.Equal(fakeTable.RenderRowCallCount(), 0)
	})
}

func TestJobTable_Sort(t *testing.T) {
	r := require.New(t)

	newData := func() []*models.Job {
		now := time.Now()
		return []*models.Job{
			{ID: "b-id", Name: "beta", Type: "batch", Namespace: "ns2", Status: "running", SubmitTime: now.Add(-1 * time.Hour)},
			{ID: "a-id", Name: "alpha", Type: "service", Namespace: "ns1", Status: "dead", SubmitTime: now.Add(-2 * time.Hour)},
			{ID: "c-id", Name: "gamma", Type: "service", Namespace: "ns3", Status: "pending", SubmitTime: now},
		}
	}

	t.Run("CycleSort advances through columns in order and wraps, always ascending", func(t *testing.T) {
		jt := component.NewJobsTable()

		expectedOrder := []string{
			component.SortByID,
			component.SortByName,
			component.SortByType,
			component.SortByNamespace,
			component.SortByStatus,
			component.SortBySubmitTime,
			component.SortByUptime,
			component.SortByID,
		}

		for _, expected := range expectedOrder {
			jt.CycleSort()
			r.Equal(expected, jt.Props.SortColumn)
			r.True(jt.Props.SortAscending)
		}
	})

	t.Run("FlipSortDirection toggles direction and is a no-op with no active sort", func(t *testing.T) {
		jt := component.NewJobsTable()

		// No-op: no column is sorted yet.
		jt.FlipSortDirection()
		r.Empty(jt.Props.SortColumn)
		r.False(jt.Props.SortAscending)

		jt.CycleSort()
		r.True(jt.Props.SortAscending)

		jt.FlipSortDirection()
		r.False(jt.Props.SortAscending)

		jt.FlipSortDirection()
		r.True(jt.Props.SortAscending)
	})

	t.Run("Render sorts rows by ID ascending then descending", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		jt := component.NewJobsTable()
		jt.Table = fakeTable
		jt.Props.Data = newData()
		jt.Props.SelectJob = func(id string) {}
		jt.Props.HandleNoResources = func(format string, args ...interface{}) {}
		jt.Bind(tview.NewFlex())

		jt.Props.SortColumn = component.SortByID
		jt.Props.SortAscending = true
		r.NoError(jt.Render())

		row1, _, _ := fakeTable.RenderRowArgsForCall(0)
		row2, _, _ := fakeTable.RenderRowArgsForCall(1)
		row3, _, _ := fakeTable.RenderRowArgsForCall(2)
		r.Equal("a-id", row1[1])
		r.Equal("b-id", row2[1])
		r.Equal("c-id", row3[1])

		header := fakeTable.RenderHeaderArgsForCall(fakeTable.RenderHeaderCallCount() - 1)
		r.Contains(header, "ID ▲")

		jt.Props.SortAscending = false
		r.NoError(jt.Render())

		row1, _, _ = fakeTable.RenderRowArgsForCall(3)
		row2, _, _ = fakeTable.RenderRowArgsForCall(4)
		row3, _, _ = fakeTable.RenderRowArgsForCall(5)
		r.Equal("c-id", row1[1])
		r.Equal("b-id", row2[1])
		r.Equal("a-id", row3[1])

		header = fakeTable.RenderHeaderArgsForCall(fakeTable.RenderHeaderCallCount() - 1)
		r.Contains(header, "ID ▼")
	})

	t.Run("Render sorts rows by SubmitTime and Uptime as independent directions", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		jt := component.NewJobsTable()
		jt.Table = fakeTable
		jt.Props.Data = newData()
		jt.Props.SelectJob = func(id string) {}
		jt.Props.HandleNoResources = func(format string, args ...interface{}) {}
		jt.Bind(tview.NewFlex())

		// SubmitTime ascending: oldest submit time (a-id) first.
		jt.Props.SortColumn = component.SortBySubmitTime
		jt.Props.SortAscending = true
		r.NoError(jt.Render())

		row1, _, _ := fakeTable.RenderRowArgsForCall(0)
		r.Equal("a-id", row1[1])

		// Uptime ascending: smallest uptime, i.e. most recently submitted
		// (c-id), first - the inverse of SubmitTime ascending.
		jt.Props.SortColumn = component.SortByUptime
		jt.Props.SortAscending = true
		r.NoError(jt.Render())

		row1, _, _ = fakeTable.RenderRowArgsForCall(3)
		r.Equal("c-id", row1[1])
	})

	t.Run("Render leaves natural order and header unchanged with no sort set", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		jt := component.NewJobsTable()
		jt.Table = fakeTable
		jt.Props.Data = newData()
		jt.Props.SelectJob = func(id string) {}
		jt.Props.HandleNoResources = func(format string, args ...interface{}) {}
		jt.Bind(tview.NewFlex())

		r.NoError(jt.Render())

		row1, _, _ := fakeTable.RenderRowArgsForCall(0)
		r.Equal("b-id", row1[1])

		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.TableHeaderJobs, header)
	})
}

func TestJobTable_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When SelectDeployment is not set", func(t *testing.T) {
		jt := component.NewJobsTable()

		jt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		jt.Bind(slot)

		// It doesn't error
		err := jt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When HandleNoResources is not set", func(t *testing.T) {
		jt := component.NewJobsTable()

		jt.Props.SelectJob = func(id string) {}

		slot := tview.NewFlex()
		jt.Bind(slot)

		// It doesn't error
		err := jt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When the component isn't bound", func(t *testing.T) {
		jt := component.NewJobsTable()

		jt.Props.SelectJob = func(id string) {}
		jt.Props.HandleNoResources = func(format string, args ...interface{}) {}

		// It doesn't error
		err := jt.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})
}
