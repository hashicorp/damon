// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package layout_test

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/layout"
)

func TestDefaultLayout(t *testing.T) {
	r := require.New(t)

	l := layout.New(layout.Default)

	r.NotNil(l.Container)
	r.IsType(l.Container, &tview.Application{})

	r.NotNil(l.Pages)
	r.IsType(l.Pages, &tview.Pages{})
	r.Equal(l.Pages.GetPageCount(), 1)
	r.True(l.Pages.HasPage("main"))

	r.NotNil(l.Header)
	r.NotNil(l.Header.SlotInfo)
	r.IsType(l.Header.SlotInfo, &tview.Flex{})

	r.NotNil(l.Header.SlotCmd)
	r.IsType(l.Header.SlotCmd, &tview.Flex{})

	r.NotNil(l.Header.SlotLogo)
	r.IsType(l.Header.SlotLogo, &tview.Flex{})

	r.NotNil(l.Elements)
	r.NotNil(l.Elements.ClusterInfo)
	r.IsType(l.Elements.ClusterInfo, &tview.Flex{})

	r.NotNil(l.Elements.Dropdowns)
	r.IsType(l.Elements.Dropdowns, &tview.Flex{})

	r.NotNil(l.Body)
	r.IsType(l.Body, &tview.Flex{})

	r.NotNil(l.Footer)
	r.IsType(l.Footer, &tview.Flex{})

	r.NotNil(l.MainPage)
	r.IsType(l.MainPage, &tview.Flex{})
}
