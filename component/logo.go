// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"strings"

	"github.com/rivo/tview"

	primitive "github.com/hcjulz/damon/primitives"
)

var LogoASCII = []string{
	`[#00b57c]    .___                             `,
	`  __| _/____    _____   ____   ____  `,
	` / __ |\__  \  /     \ /  _ \ /    \ `,
	`/ /_/ | / __ \|  Y Y  (  <_> )   |  \`,
	`\____ |(____  /__|_|  /\____/|___|  /`,
	`     \/     \/      \/            \/ `,
	`[#26ffe6]HashiCorp Nomad - Terminal Dashboard`,
}

type Logo struct {
	TextView TextView
	slot     *tview.Flex
}

func NewLogo() *Logo {
	t := primitive.NewTextView(tview.AlignRight)
	return &Logo{
		TextView: t,
	}
}

func (l *Logo) Render() error {
	if l.slot == nil {
		return ErrComponentNotBound
	}

	logo := strings.Join(LogoASCII, "\n")

	l.TextView.SetText(logo)
	l.slot.AddItem(l.TextView.Primitive(), 0, 1, false)
	return nil
}

func (l *Logo) Bind(slot *tview.Flex) {
	l.slot = slot
}
