package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	primitive "github.com/hcjulz/damon/primitives"
)

type GenericModal struct {
	Modal Modal
	Props *ModalProps
	pages *tview.Pages
}

type ModalProps struct {
	ID   string
	Done DoneModalFunc
}

func NewModal(id, title string, buttons []string, c tcell.Color) *GenericModal {
	modal := primitive.NewModal(title, buttons, c)

	return &GenericModal{
		Modal: modal,
		Props: &ModalProps{ID: id},
	}
}

func (m *GenericModal) Render(msg string) error {
	if m.Props.Done == nil {
		return ErrComponentPropsNotSet
	}

	if m.pages == nil {
		return ErrComponentNotBound
	}

	m.Modal.SetFocus(0)
	m.Modal.SetDoneFunc(m.Props.Done)
	m.Modal.SetText(msg)
	m.pages.AddPage(m.Props.ID, m.Modal.Container(), true, true)

	return nil
}

func (m *GenericModal) Bind(pages *tview.Pages) {
	m.pages = pages
}
