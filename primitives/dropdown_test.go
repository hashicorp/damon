package primitives_test

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

func TestDropDown(t *testing.T) {
	r := require.New(t)

	dd := primitives.NewDropDown("test")
	p := dd.Primitive().(*tview.DropDown)

	r.Equal(p.GetBackgroundColor(), styles.TcellBackgroundColor)
	r.Equal(p.GetLabel(), "test")
}

func TestDropDown_Options(t *testing.T) {
	r := require.New(t)

	dd := primitives.NewDropDown("test")
	p := dd.Primitive().(*tview.DropDown)

	dd.SetOptions([]string{"opt", "opt2"}, nil)

	dd.SetCurrentOption(0)
	index, text := p.GetCurrentOption()
	r.Equal(text, "opt")
	r.Equal(index, 0)

	dd.SetCurrentOption(1)
	index, text = p.GetCurrentOption()
	r.Equal(text, "opt2")
	r.Equal(index, 1)
}

func TestDropDown_SetSelectedFunc(t *testing.T) {
	r := require.New(t)

	dd := primitives.NewDropDown("test")
	p := dd.Primitive().(*tview.DropDown)

	dd.SetOptions([]string{"opt", "opt2"}, nil)

	var selected bool
	dd.SetSelectedFunc(func(text string, index int) {
		selected = true
	})

	dd.SetCurrentOption(0)
	index, text := p.GetCurrentOption()
	r.Equal(text, "opt")
	r.Equal(index, 0)
	r.True(selected)
}
