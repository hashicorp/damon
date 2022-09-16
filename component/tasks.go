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
	TableTitleTasks = "Tasks"
)

var (
	TableHeaderTasks = []string{
		LabelName,
		LabelState,
		LabelDriver,
		LabelImage,
		LabelLastEvent,
	}
)

type SelectTaskFunc func(allocID, taskID string)

type TaskTable struct {
	Table Table
	Props *TaskTableProps

	slot        *tview.Flex
	keyBindings map[tcell.Key]func(event *tcell.EventKey)
}

type TaskTableProps struct {
	SelectTask        SelectTaskFunc
	HandleNoResources models.HandlerFunc

	AllocationID string

	Data []*models.Task
}

func NewTaskTable() *TaskTable {
	t := primitive.NewTable()

	tt := &TaskTable{
		Table:       t,
		Props:       &TaskTableProps{},
		keyBindings: map[tcell.Key]func(event *tcell.EventKey){},
	}

	tt.Table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if fn, ok := tt.keyBindings[event.Key()]; ok {
			fn(event)
		}

		return event
	})

	return tt
}

func (t *TaskTable) Bind(slot *tview.Flex) {
	t.slot = slot
}

func (t *TaskTable) Render() error {
	if t.Props.SelectTask == nil {
		return ErrComponentPropsNotSet
	}

	if t.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	if t.slot == nil {
		return ErrComponentNotBound
	}

	t.reset()

	if len(t.Props.Data) == 0 {
		t.Props.HandleNoResources(
			"%sno tasks available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	t.Table.SetSelectedFunc(t.taskSelected)

	t.Table.RenderHeader(TableHeaderTasks)
	t.renderRows()

	t.Table.SetTitle(fmt.Sprintf("%s (Allocation: %s)", TableTitleTasks, t.Props.AllocationID))
	t.slot.AddItem(t.Table.Primitive(), 0, 1, false)

	return nil
}

func (t *TaskTable) reset() {
	t.slot.Clear()
	t.Table.Clear()
}

func (t *TaskTable) renderRows() {
	for i, task := range t.Props.Data {
		row := []string{
			task.Name,
			task.State,
			task.Driver,
		}

		if image, ok := task.Config["image"]; ok {
			row = append(row, image.(string))
		}

		row = append(row, task.Events[len(task.Events)-1].DisplayMessage)
		// row = append(row, strconv.Itoa(task.CPU))
		// row = append(row, strconv.Itoa(task.MemoryMB))

		index := i + 1

		c := t.getCellColor(task.State)
		t.Table.RenderRow(row, index, c)
	}
}

func (t *TaskTable) getCellColor(status string) tcell.Color {
	c := tcell.ColorWhite

	switch status {
	case models.StatusDead:
		c = tcell.ColorGray
	case models.StatusFailed:
		c = tcell.ColorRed
	case models.StatusPending:
		c = tcell.ColorYellow
	}

	return c
}

func (t *TaskTable) taskSelected(row, column int) {
	taskName := t.Table.GetCellContent(row, 0)
	t.Props.SelectTask(taskName, t.Props.AllocationID)
}

func (t *TaskTable) GetNameForSelection() string {
	row, _ := t.Table.GetSelection()
	return t.Table.GetCellContent(row, 0)
}

func (t *TaskTable) BindKey(key tcell.Key, fn func(event *tcell.EventKey)) {
	t.keyBindings[key] = fn
}
