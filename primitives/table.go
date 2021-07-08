package primitives

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/styles"
)

// Table is a wrapper of a tview.Table primitive.
// It applies the damon look to the tviw.Table.
type Table struct {
	primitive *tview.Table
	color     tcell.Color
}

func NewTable() *Table {
	t := tview.NewTable()
	t.SetBorder(true)
	t.SetTitleColor(styles.TcellColorHighlighPrimary)
	t.SetSelectable(true, false)
	t.SetFixed(1, 1)
	t.SetBorderPadding(0, 0, 1, 1)
	t.SetBorderColor(styles.TcellColorStandard)

	return &Table{
		primitive: t,
	}
}

func (t *Table) RenderHeader(data []string) {
	for i, h := range data {
		c := tcell.GetColor(styles.StandardColorHex)
		t.primitive.SetCell(0, i, tview.NewTableCell(h).
			SetTextColor(c).
			SetSelectable(false),
		)
	}
}

func (t *Table) SetTitle(format string, args ...interface{}) {
	t.primitive.SetTitle(fmt.Sprintf(format, args...))
}

func (t *Table) GetCellContent(row, column int) string {
	cell := t.primitive.GetCell(row, column)
	return cell.Text
}

func (t *Table) GetSelection() (row, column int) {
	return t.primitive.GetSelection()
}

func (t *Table) Clear() {
	t.primitive.Clear()
}

func (t *Table) RenderRow(data []string, index int, c tcell.Color) {
	for i, r := range data {
		t.primitive.SetCell(index, i,
			tview.NewTableCell(r).SetTextColor(c).SetExpansion(1),
		)
	}
}

func (t *Table) SetSelectedFunc(fn func(row, column int)) {
	t.primitive.SetSelectedFunc(fn)
}

func (t *Table) Primitive() tview.Primitive {
	return t.primitive
}
