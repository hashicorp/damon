// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"github.com/rivo/tview"

	primitive "github.com/hcjulz/damon/primitives"
)

type ClusterInfo struct {
	TextView TextView
	Props    *ClusterInfoProps

	slot *tview.Flex
}

type ClusterInfoProps struct {
	Info string
}

func NewClusterInfo() *ClusterInfo {
	return &ClusterInfo{
		TextView: primitive.NewTextView(tview.AlignLeft),
		Props:    &ClusterInfoProps{},
	}
}

func (c *ClusterInfo) Render() error {
	if c.slot == nil {
		return ErrComponentNotBound
	}

	c.TextView.SetText(c.Props.Info)
	c.slot.AddItem(c.TextView.Primitive(), 0, 1, false)

	return nil
}

func (c *ClusterInfo) Bind(slot *tview.Flex) {
	c.slot = slot
}
