package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SelectorModal struct {
	Modal Selector
	Props *SelectorProps
	pages *tview.Pages
}

type SelectorProps struct {
	ID    string
	Items []string
}

func NewSelectorModal() *SelectorModal {
	return &SelectorModal{}
}

func (s *SelectorModal) Render() error {
	if s.pages == nil {
		return ErrComponentNotBound
	}

	if s.Props.Items == nil {
		return ErrComponentPropsNotSet
	}

	table := s.Modal.GetTable()
	table.Clear()

	// SetSelectedFunc

	for i, v := range s.Props.Items {
		table.RenderRow([]string{v}, i, tcell.ColorWhite)
	}

	s.pages.AddPage(s.Props.ID, s.Modal.Container(), true, true)

	return nil
}

func (s *SelectorModal) Bind(pages *tview.Pages) {
	s.pages = pages
}
