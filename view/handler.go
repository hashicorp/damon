package view

import (
	"fmt"

	"github.com/rivo/tview"
)

func (v *View) handleNoResources(text string, args ...interface{}) {
	msg := fmt.Sprintf(text, args...)
	info := tview.
		NewTextView().
		SetDynamicColors(true).
		SetText(msg).
		SetTextAlign(tview.AlignCenter)
	info.SetBorder(true)

	v.Layout.Body.AddItem(info, 0, 1, false)
}

func (v *View) err(err error, msg string) {
	if err != nil {
		fmt.Sprintf("%s: %s", msg, err.Error())
		v.handleError(msg)
	}
}

func (v *View) handleError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	v.components.Failure.Render(msg)
	v.Layout.Container.SetFocus(v.components.Failure.Modal.Primitive())
}

func (v *View) handleInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	v.components.Info.Render(msg)
	v.Layout.Container.SetFocus(v.components.Info.Modal.Primitive())
}

func (v *View) handleFatal(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	v.components.Error.Render(msg)
	v.Layout.Container.SetFocus(v.components.Error.Modal.Primitive())
}
