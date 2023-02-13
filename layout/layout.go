// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package layout

import (
	"github.com/rivo/tview"
)

const NameMainPage = "main"
const NameErrorPage = "error"

type Layout struct {
	Container *tview.Application

	Pages    *tview.Pages
	MainPage *tview.Flex

	Header *Header
	Body   *tview.Flex
	Footer *tview.Flex

	Elements *Elements
}

type Elements struct {
	ClusterInfo *tview.Flex
	Dropdowns   *tview.Flex
}

type Header struct {
	SlotInfo *tview.Flex
	SlotCmd  *tview.Flex
	SlotLogo *tview.Flex
}

func EnableMouse(l *Layout) {
	l.Container.EnableMouse(true)
}

func New(options ...func(*Layout)) *Layout {
	v := &Layout{}

	for _, opt := range options {
		opt(v)
	}

	return v
}

func Default(l *Layout) {
	l.Header = &Header{}
	l.Elements = &Elements{}

	l.Elements.ClusterInfo = tview.NewFlex()
	l.Elements.Dropdowns = tview.NewFlex().SetDirection(tview.FlexRow)

	l.Header.SlotInfo = tview.NewFlex().SetDirection(tview.FlexRow)
	l.Header.SlotInfo.AddItem(l.Elements.ClusterInfo, 0, 1, false)
	l.Header.SlotInfo.AddItem(l.Elements.Dropdowns, 0, 1, false)

	l.Header.SlotCmd = tview.NewFlex()
	l.Header.SlotLogo = tview.NewFlex()

	header := tview.NewFlex().
		AddItem(l.Header.SlotInfo, 0, 1, false).
		AddItem(l.Header.SlotCmd, 0, 1, false).
		AddItem(l.Header.SlotLogo, 0, 1, false)

	header.SetBorderPadding(1, 1, 2, 2)

	footer := tview.NewFlex()
	body := tview.NewFlex()

	mainPage := tview.NewFlex().SetDirection(tview.FlexRow)
	mainPage.
		AddItem(header, 0, 4, false).
		AddItem(body, 0, 12, false).
		AddItem(footer, 0, 0, false)

	pages := tview.NewPages()
	pages.AddPage(NameMainPage, mainPage, true, true)

	l.Body = body
	l.Footer = footer

	l.MainPage = mainPage
	l.Pages = pages

	l.Container = tview.NewApplication().
		SetRoot(pages, true).
		SetFocus(pages)

}
