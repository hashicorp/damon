package primitives_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

func TestTable(t *testing.T) {
	tview.Styles.PrimitiveBackgroundColor = styles.TcellBackgroundColor

	r := require.New(t)

	tb := primitives.NewTable()
	p := tb.Primitive().(*tview.Table)

	r.Equal(p.GetBackgroundColor(), styles.TcellBackgroundColor)
}

func TestTable_Render(t *testing.T) {
	r := require.New(t)

	tb := primitives.NewTable()

	header := []string{"col1", "col2", "col3"}
	row1 := []string{"row1-1", "row1-2", "row1-3"}
	row2 := []string{"row2-1", "row2-2", "row2-3"}

	tb.RenderHeader(header)
	tb.RenderRow(row1, 1, tcell.ColorWhite)
	tb.RenderRow(row2, 2, tcell.ColorWhite)

	h1 := tb.GetCellContent(0, 0)
	h2 := tb.GetCellContent(0, 1)
	h3 := tb.GetCellContent(0, 2)

	r11 := tb.GetCellContent(1, 0)
	r12 := tb.GetCellContent(1, 1)
	r13 := tb.GetCellContent(1, 2)

	r21 := tb.GetCellContent(2, 0)
	r22 := tb.GetCellContent(2, 1)
	r23 := tb.GetCellContent(2, 2)

	r.Equal(h1, "col1")
	r.Equal(h2, "col2")
	r.Equal(h3, "col3")

	r.Equal(r11, "row1-1")
	r.Equal(r12, "row1-2")
	r.Equal(r13, "row1-3")

	r.Equal(r21, "row2-1")
	r.Equal(r22, "row2-2")
	r.Equal(r23, "row2-3")

}

func TestTable_Clear(t *testing.T) {
	r := require.New(t)

	tb := primitives.NewTable()
	p := tb.Primitive().(*tview.Table)

	header := []string{"col1", "col2", "col3"}
	row1 := []string{"row1-1", "row1-2", "row1-3"}
	row2 := []string{"row2-1", "row2-2", "row2-3"}

	tb.RenderHeader(header)
	tb.RenderRow(row1, 1, tcell.ColorWhite)
	tb.RenderRow(row2, 2, tcell.ColorWhite)

	tb.Clear()

	r.Equal(p.GetColumnCount(), 0)
	r.Equal(p.GetRowCount(), 0)
}

func TestTable_GetSelection(t *testing.T) {
	r := require.New(t)

	tb := primitives.NewTable()
	p := tb.Primitive().(*tview.Table)

	header := []string{"col1", "col2", "col3"}
	row1 := []string{"row1-1", "row1-2", "row1-3"}
	row2 := []string{"row2-1", "row2-2", "row2-3"}

	tb.RenderHeader(header)
	tb.RenderRow(row1, 1, tcell.ColorWhite)
	tb.RenderRow(row2, 2, tcell.ColorWhite)

	p.Select(2, 2)

	row, col := tb.GetSelection()

	r.Equal(row, 2)
	r.Equal(col, 2)
}

func TestTable_SetTitle(t *testing.T) {
	r := require.New(t)

	tb := primitives.NewTable()
	p := tb.Primitive().(*tview.Table)

	tb.SetTitle("test")

	r.Equal(p.GetTitle(), "test")
}
