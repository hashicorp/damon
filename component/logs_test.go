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
	"github.com/hcjulz/damon/styles"
)

func TestLogs_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is data to render", func(t *testing.T) {
		textView := &componentfakes.FakeTextView{}
		logs := component.NewLogger()
		logs.TextView = textView
		logs.Props.HandleNoResources = func(format string, args ...interface{}) {}
		logs.Props.Data = []byte("logs")

		logs.Bind(tview.NewFlex())

		err := logs.Render()
		r.NoError(err)

		text := textView.SetTextArgsForCall(0)
		r.Equal(string(text), "logs")
	})

	t.Run("When there is no data to render", func(t *testing.T) {
		textView := &componentfakes.FakeTextView{}
		logs := component.NewLogger()
		logs.TextView = textView

		var handleNoResourcesCalled bool
		logs.Props.HandleNoResources = func(format string, args ...interface{}) {
			handleNoResourcesCalled = true

			r.Equal("%sWHOOOPS, no Logs found", format)
			r.Len(args, 1)
			r.Equal(args[0], styles.HighlightSecondaryTag)
		}

		logs.Bind(tview.NewFlex())

		err := logs.Render()
		r.NoError(err)

		r.True(handleNoResourcesCalled)
	})
}

func TestLogs_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the component is not bound", func(t *testing.T) {
		logs := component.NewLogger()
		logs.Props.HandleNoResources = func(format string, args ...interface{}) {}

		err := logs.Render()
		r.Error(err)

		r.True(errors.Is(err, component.ErrComponentNotBound))
		r.EqualError(err, "component not bound")
	})

	t.Run("When the component props are not set", func(t *testing.T) {
		logs := component.NewLogger()
		logs.Bind(tview.NewFlex())

		err := logs.Render()
		r.Error(err)

		r.True(errors.Is(err, component.ErrComponentPropsNotSet))
		r.EqualError(err, "component properties not set")
	})
}
