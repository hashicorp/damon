package component_test

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
)

func TestJobStatus_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("When there is data to render", func(t *testing.T) {
		textView := &componentfakes.FakeTextView{}
		jobStatus := component.NewJobStatus()
		jobStatus.TextView = textView
		jobStatus.Status = "Sample Status"

		jobStatus.Bind(tview.NewFlex())

		err := jobStatus.Render()
		r.NoError(err)

		text := textView.SetTextArgsForCall(0)
		r.Equal(text, "Sample Status")
	})

	t.Run("When there is no data to render", func(t *testing.T) {
		textView := &componentfakes.FakeTextView{}
		jobStatus := component.NewJobStatus()
		jobStatus.TextView = textView
		jobStatus.Status = ""

		jobStatus.Bind(tview.NewFlex())

		err := jobStatus.Render()
		r.NoError(err)

		text := textView.SetTextArgsForCall(0)
		r.Equal(text, "Status not available.")
	})
}

func TestJobStatus_Sad(t *testing.T) {
	r := require.New(t)

	t.Run("When the component is not bound", func(t *testing.T) {
		jobStatus := component.NewJobStatus()

		err := jobStatus.Render()
		r.Error(err)

		r.True(errors.Is(err, component.ErrComponentNotBound))
		r.EqualError(err, "component not bound")
	})
}