// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package primitives

import (
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/styles"
)

type DropDown struct {
	primitive *tview.DropDown
}

func NewDropDown(label string) *DropDown {
	dd := tview.NewDropDown()
	dd.SetLabel(label)
	dd.SetBackgroundColor(styles.TcellBackgroundColor)
	dd.SetCurrentOption(0)
	dd.SetFieldBackgroundColor(styles.TcellBackgroundColor)
	dd.SetFieldTextColor(styles.TcellColorStandard)

	return &DropDown{dd}
}

func (d *DropDown) SetOptions(options []string, selected func(text string, index int)) {
	d.primitive.SetOptions(options, selected)
}

func (d *DropDown) SetCurrentOption(index int) {
	d.primitive.SetCurrentOption(index)
}

func (d *DropDown) SetSelectedFunc(selected func(text string, index int)) {
	d.primitive.SetSelectedFunc(selected)
}

func (d *DropDown) Primitive() tview.Primitive {
	return d.primitive
}
