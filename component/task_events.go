package component

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/hashicorp/nomad/api"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

const (
	TableTitleTaskEvents = "TaskEvents"
)

var (
	TableHeaderTaskEvents = []string{
		LabelTime,
		LabelType,
		LabelMessage,
	}
)

type TaskEventsTable struct {
	Table Table
	Props *TaskEventsTableProps

	slot *tview.Flex
}

type TaskEventsTableProps struct {
	HandleNoResources models.HandlerFunc
	Data              []*api.TaskEvent
	AllocID           string
}

func NewTaskEventsTable() *TaskEventsTable {
	return &TaskEventsTable{
		Table: primitive.NewTable(),
		Props: &TaskEventsTableProps{},
	}
}

func (t *TaskEventsTable) Bind(slot *tview.Flex) {
	slot.SetTitle("Events")
	t.slot = slot
}

func (t *TaskEventsTable) Render() error {
	if t.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	if t.slot == nil {
		return ErrComponentNotBound
	}

	t.slot.Clear()
	t.Table.Clear()

	if len(t.Props.Data) == 0 {
		t.Props.HandleNoResources(
			"%sno task events available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	t.Table.SetTitle(fmt.Sprintf("%s (%s)", TableTitleTaskEvents, t.Props.AllocID))

	t.Table.RenderHeader(TableHeaderTaskEvents)
	t.renderRows()

	t.slot.AddItem(t.Table.Primitive(), 0, 1, false)
	return nil
}

func (t *TaskEventsTable) renderRows() {
	for i, e := range t.Props.Data {
		row := []string{
			time.Unix(0, e.Time).Format(time.RFC3339),
			e.Type,
			e.DisplayMessage,
		}

		index := i + 1

		t.Table.RenderRow(row, index, tcell.ColorWhite)
	}
}

func (t *TaskEventsTable) GetIDForSelection() string {
	row, _ := t.Table.GetSelection()
	return t.Table.GetCellContent(row, 0)
}
