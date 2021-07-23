package component_test

import (
	"errors"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
)

func TestClusterInfo_Happy(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	clusterInfo := component.NewClusterInfo()
	clusterInfo.TextView = textView
	clusterInfo.Props.Info = "info"

	clusterInfo.Bind(tview.NewFlex())

	err := clusterInfo.Render()
	r.NoError(err)

	text := textView.SetTextArgsForCall(0)
	r.Equal(text, "info")
}

func TestClusterInfo_Render_Sad(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	clusterInfo := component.NewClusterInfo()
	clusterInfo.TextView = textView
	clusterInfo.Props.Info = "info"

	err := clusterInfo.Render()
	r.Error(err)

	r.True(errors.Is(err, component.ErrComponentNotBound))
	r.EqualError(err, "component not bound")
}
