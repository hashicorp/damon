// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/primitives"
)

const pageNameSelector = "selector"

type SelectorModal struct {
	Modal       Selector
	Props       *SelectorProps
	pages       *tview.Pages
	keyBindings map[tcell.Key]func()
}

type SelectorProps struct {
	Items        []string
	AllocationID string
}

func NewSelectorModal() *SelectorModal {
	s := &SelectorModal{
		Modal:       primitives.NewSelectionModal(),
		Props:       &SelectorProps{},
		keyBindings: map[tcell.Key]func(){},
	}

	s.Modal.GetTable().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if fn, ok := s.keyBindings[event.Key()]; ok {
			fn()
		}

		return event
	})

	return s
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

	for i, v := range s.Props.Items {
		table.RenderRow([]string{v}, i, tcell.ColorWhite)
	}

	s.Modal.GetTable().SetTitle("Select a Task (alloc: %s)", s.Props.AllocationID)

	s.pages.AddPage(pageNameSelector, s.Modal.Container(), true, true)

	return nil
}

func (s *SelectorModal) Bind(pages *tview.Pages) {
	s.pages = pages
}

func (s *SelectorModal) SetSelectedFunc(fn func(task string)) {
	s.Modal.GetTable().SetSelectedFunc(func(row, column int) {
		task := s.Modal.GetTable().GetCellContent(row, 0)
		fn(task)
		s.Close()
	})
}

func (s *SelectorModal) Close() {
	s.pages.RemovePage(pageNameSelector)
}

func (s *SelectorModal) BindKey(key tcell.Key, fn func()) {
	s.keyBindings[key] = fn
}
