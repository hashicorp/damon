// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"github.com/rivo/tview"

	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

const PageNameInfo = "info"

type Info struct {
	Modal Modal
	Props *InfoProps
	pages *tview.Pages
}

type InfoProps struct {
	Done DoneModalFunc
}

func NewInfo() *Info {
	buttons := []string{"OK"}
	modal := primitive.NewModal("Info", buttons, styles.TcellColorModalInfo)

	return &Info{
		Modal: modal,
		Props: &InfoProps{},
	}
}

func (i *Info) Render(msg string) error {
	if i.Props.Done == nil {
		return ErrComponentPropsNotSet
	}

	if i.pages == nil {
		return ErrComponentNotBound
	}

	i.Modal.SetDoneFunc(i.Props.Done)
	i.Modal.SetText(msg)
	i.pages.AddPage(PageNameInfo, i.Modal.Container(), true, true)

	return nil
}

func (i *Info) Bind(pages *tview.Pages) {
	i.pages = pages
}
