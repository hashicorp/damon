// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/primitives"
)

const gotoLinePlaceholder = "(line number, hit enter or esc to leave)"

type GotoLine struct {
	InputField InputField
	Props      *GotoLineProps
	slot       *tview.Flex
}

type GotoLineProps struct {
	DoneFunc SetDoneFunc
}

func NewGotoLine() *GotoLine {
	gl := &GotoLine{}
	gl.Props = &GotoLineProps{}
	gl.InputField = primitives.NewInputField(": ", gotoLinePlaceholder)

	return gl
}

func (gl *GotoLine) Render() error {
	if err := gl.validate(); err != nil {
		return err
	}

	gl.InputField.SetDoneFunc(gl.Props.DoneFunc)
	gl.slot.AddItem(gl.InputField.Primitive(), 0, 2, false)
	return nil
}

func (gl *GotoLine) validate() error {
	if gl.Props.DoneFunc == nil {
		return ErrComponentPropsNotSet
	}

	if gl.slot == nil {
		return ErrComponentNotBound
	}

	return nil
}

func (gl *GotoLine) Bind(slot *tview.Flex) {
	gl.slot = slot
}
