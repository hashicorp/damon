// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

const (
	TableTitleTaskGroups = "TaskGroups"
)

var (
	TableHeaderTaskGroups = []string{
		LabelName,
		LabelJobID,
		LabelStarting,
		LabelQueued,
		LabelRunning,
		LabelComplete,
		LabelFailed,
		LabelLost,
	}
)

type SelectTaskGroupFunc func(ID string)

type TaskGroupTable struct {
	Table Table
	Props *TaskGroupTableProps

	slot *tview.Flex
}

type TaskGroupTableProps struct {
	SelectTaskGroup   SelectTaskGroupFunc
	HandleNoResources models.HandlerFunc
	Data              []*models.TaskGroup
	JobID             string
}

func NewTaskGroupTable() *TaskGroupTable {
	t := primitive.NewTable()

	return &TaskGroupTable{
		Table: t,
		Props: &TaskGroupTableProps{},
	}
}

func (t *TaskGroupTable) Bind(slot *tview.Flex) {
	t.slot = slot
}

func (t *TaskGroupTable) Render() error {
	if t.Props.SelectTaskGroup == nil || t.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	if t.slot == nil {
		return ErrComponentNotBound
	}

	t.slot.Clear()
	t.Table.Clear()

	if len(t.Props.Data) == 0 {
		t.Props.HandleNoResources(
			"%sno TaskGroups available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	t.Table.SetSelectedFunc(t.taskGroupSelected)
	t.Table.SetTitle(fmt.Sprintf("%s (%s)", TableTitleTaskGroups, t.Props.JobID))

	t.Table.RenderHeader(TableHeaderTaskGroups)
	t.renderRows()

	t.slot.AddItem(t.Table.Primitive(), 0, 1, false)
	return nil
}

func (t *TaskGroupTable) renderRows() {
	for i, tg := range t.Props.Data {
		row := []string{
			tg.Name,
			tg.JobID,
			fmt.Sprint(tg.Starting),
			fmt.Sprint(tg.Queued),
			fmt.Sprint(tg.Running),
			fmt.Sprint(tg.Complete),
			fmt.Sprint(tg.Failed),
			fmt.Sprint(tg.Lost),
		}

		index := i + 1

		t.Table.RenderRow(row, index, tcell.ColorWhite)
	}
}

func (t *TaskGroupTable) taskGroupSelected(row, column int) {
	jobID := t.Table.GetCellContent(row, 0)
	t.Props.SelectTaskGroup(jobID)
}
