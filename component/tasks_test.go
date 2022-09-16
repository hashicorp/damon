package component_test

import (
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/styles"
)

func TestTasks_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		at := component.NewTaskTable()

		at.Table = fakeTable
		at.Props.AllocationID = "japan"
		at.Props.Data = []*models.Task{
			{
				Name:     "task-1",
				Driver:   "docker",
				State:    "running",
				CPU:      100,
				MemoryMB: 10,
				DiskMB:   1000,
				Config: map[string]interface{}{
					"image": "docker-image",
				},
				Env: map[string]string{
					"env-key": "env-value",
				},
				Events: []*api.TaskEvent{
					{
						DisplayMessage: "msg",
						Type:           "type",
					},
				},
			},
			{
				Name:     "task-2",
				Driver:   "docker",
				State:    "failed",
				CPU:      100,
				MemoryMB: 10,
				DiskMB:   1000,
				Config: map[string]interface{}{
					"image": "docker-image",
				},
				Env: map[string]string{
					"env-key": "env-value",
				},
				Events: []*api.TaskEvent{
					{
						DisplayMessage: "msg",
						Type:           "type",
					},
				},
			},
			{
				Name:     "task-3",
				Driver:   "docker",
				State:    "dead",
				CPU:      100,
				MemoryMB: 10,
				DiskMB:   1000,
				Config: map[string]interface{}{
					"image": "docker-image",
				},
				Env: map[string]string{
					"env-key": "env-value",
				},
				Events: []*api.TaskEvent{
					{
						DisplayMessage: "msg",
						Type:           "type",
					},
				},
			},
			{
				Name:     "task-4",
				Driver:   "docker",
				State:    "pending",
				CPU:      100,
				MemoryMB: 10,
				DiskMB:   1000,
				Config: map[string]interface{}{
					"image": "docker-image",
				},
				Env: map[string]string{
					"env-key": "env-value",
				},
				Events: []*api.TaskEvent{
					{
						DisplayMessage: "msg",
						Type:           "type",
					},
				},
			},
		}

		var selectCalled bool
		at.Props.SelectTask = func(allocID, taskName string) {
			selectCalled = true
		}

		at.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		at.Bind(slot)

		// It doesn't error
		err := at.Render()
		r.NoError(err)

		// It sets the correct selected function
		selectedFunc := fakeTable.SetSelectedFuncArgsForCall(0)
		selectedFunc(1, 1)
		r.True(selectCalled)

		// It renders the correct number of header rows
		renderHeaderCount := fakeTable.RenderHeaderCallCount()
		r.Equal(renderHeaderCount, 1)

		// It renders the correct header values
		header := fakeTable.RenderHeaderArgsForCall(0)
		r.Equal(component.TableHeaderTasks, header)

		// It renders the correct number of rows
		renderRowCallCount := fakeTable.RenderRowCallCount()
		r.Equal(renderRowCallCount, 4)

		row1, index1, c1 := fakeTable.RenderRowArgsForCall(0)
		row2, index2, c2 := fakeTable.RenderRowArgsForCall(1)
		row3, index3, c3 := fakeTable.RenderRowArgsForCall(2)
		row4, index4, c4 := fakeTable.RenderRowArgsForCall(3)

		expectedRow1 := []string{"task-1", "running", "docker", "docker-image", "msg"}
		expectedRow2 := []string{"task-2", "failed", "docker", "docker-image", "msg"}
		expectedRow3 := []string{"task-3", "dead", "docker", "docker-image", "msg"}
		expectedRow4 := []string{"task-4", "pending", "docker", "docker-image", "msg"}

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
		r.Equal(c1, tcell.ColorWhite)
		r.Equal(c2, tcell.ColorRed)
		r.Equal(c3, tcell.ColorGray)
		r.Equal(c4, tcell.ColorYellow)
	})

	t.Run("When render called again", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		at := component.NewTaskTable()

		at.Table = fakeTable
		at.Props.AllocationID = "japan"
		at.Props.Data = []*models.Task{
			{
				Name:     "task-1",
				Driver:   "docker",
				State:    "running",
				CPU:      100,
				MemoryMB: 10,
				DiskMB:   1000,
				Config: map[string]interface{}{
					"image": "docker-image",
				},
				Env: map[string]string{
					"env-key": "env-value",
				},
				Events: []*api.TaskEvent{
					{
						DisplayMessage: "msg",
						Type:           "type",
					},
				},
			},
		}

		at.Props.SelectTask = func(allocID, taskName string) {}
		at.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		at.Bind(slot)

		// It doesn't error
		err := at.Render()
		r.NoError(err)

		err = at.Render()
		r.NoError(err)

		// It clears the table on each call
		r.Equal(fakeTable.ClearCallCount(), 2)
	})

	t.Run("When there is no data to render", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		at := component.NewTaskTable()

		at.Table = fakeTable
		at.Props.AllocationID = "japan"
		at.Props.Data = []*models.Task{}
		at.Props.SelectTask = func(allocID, taskName string) {}

		var handleNoResourcesCalled bool
		at.Props.HandleNoResources = func(format string, args ...interface{}) {
			handleNoResourcesCalled = true
			r.Equal("%sno tasks available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯", format)
			r.Len(args, 2)
			r.Equal(args[0], styles.HighlightPrimaryTag)
			r.Equal(args[1], styles.HighlightSecondaryTag)
		}

		slot := tview.NewFlex()
		at.Bind(slot)

		// It doesn't error
		err := at.Render()
		r.NoError(err)

		// It handled the case that there are no resources
		r.True(handleNoResourcesCalled)

		// It didn't returned after handling no resources
		r.Equal(fakeTable.RenderHeaderCallCount(), 0)
		r.Equal(fakeTable.RenderRowCallCount(), 0)
	})
}

func TestTasks_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When SelectTask is not set", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		at := component.NewTaskTable()

		at.Table = fakeTable
		at.Props.HandleNoResources = func(format string, args ...interface{}) {}

		slot := tview.NewFlex()
		at.Bind(slot)

		// It errors
		err := at.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When HandleNoResources is not set", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		at := component.NewTaskTable()

		at.Table = fakeTable
		at.Props.SelectTask = func(allocID, taskName string) {}

		slot := tview.NewFlex()
		at.Bind(slot)

		// It errors
		err := at.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component properties not set")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
	})

	t.Run("When the component isn't bound", func(t *testing.T) {
		fakeTable := &componentfakes.FakeTable{}
		at := component.NewTaskTable()

		at.Table = fakeTable
		at.Props.SelectTask = func(allocID, taskName string) {}
		at.Props.HandleNoResources = func(format string, args ...interface{}) {}

		// It errors
		err := at.Render()
		r.Error(err)

		// It provides the correct error message
		r.EqualError(err, "component not bound")

		// It is the correct error
		r.True(errors.Is(err, component.ErrComponentNotBound))
	})
}
