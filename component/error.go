package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	primitive "github.com/hcjulz/damon/primitives"
)

const PageNameError = "error"

type Error struct {
	Modal Modal
	Props *ErrorProps
	pages *tview.Pages
}

type ErrorProps struct {
	Done DoneModalFunc
}

func NewError() *Error {
	buttons := []string{"Quit", "OK"}
	modal := primitive.NewModal("Error", buttons, tcell.ColorDarkRed)

	return &Error{
		Modal: modal,
		Props: &ErrorProps{},
	}
}

func (e *Error) Render(msg string) error {
	if e.Props.Done == nil {
		return ErrComponentPropsNotSet
	}

	if e.pages == nil {
		return ErrComponentNotBound
	}

	e.Modal.SetDoneFunc(e.Props.Done)
	e.Modal.SetText(msg)
	e.pages.AddPage(PageNameError, e.Modal.Container(), true, true)
	return nil
}

func (e *Error) Bind(pages *tview.Pages) {
	e.pages = pages
}
