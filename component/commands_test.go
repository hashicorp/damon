// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
)

func TestCommands_Happy(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	cmds := component.NewCommands()
	cmds.TextView = textView
	cmds.Props.MainCommands = []string{"command1", "command2"}
	cmds.Props.ViewCommands = []string{"subCmd1", "subCmd2"}

	cmds.Bind(tview.NewFlex())

	err := cmds.Render()
	r.NoError(err)

	text := textView.SetTextArgsForCall(0)
	r.Equal(text, "command1\ncommand2\nsubCmd1\nsubCmd2")
}

func TestCommands_Sad(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	cmds := component.NewCommands()
	cmds.TextView = textView

	err := cmds.Render()
	r.Error(err)

	r.True(errors.Is(err, component.ErrComponentNotBound))
	r.EqualError(err, "component not bound")
}
