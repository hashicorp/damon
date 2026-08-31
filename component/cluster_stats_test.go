// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
	"github.com/hcjulz/damon/models"
)

func TestClusterStats_Happy(t *testing.T) {
	r := require.New(t)

	textView := &componentfakes.FakeTextView{}
	cs := component.NewClusterStats()
	cs.TextView = textView
	cs.Bind(tview.NewFlex())

	cs.Props.Loaded = true
	cs.Props.Stats = &models.ClusterStats{
		Leader: "10.0.0.1:4647",
		Nodes:  models.NodeStats{Total: 3},
	}

	err := cs.Render()
	r.NoError(err)

	text := textView.SetTextArgsForCall(0)
	r.Contains(text, "10.0.0.1:4647")
	r.Contains(text, "Nodes:   3")
}

func TestClusterStats_Sad(t *testing.T) {
	r := require.New(t)
	cs := component.NewClusterStats()

	err := cs.Render()
	r.Error(err)
	r.True(errors.Is(err, component.ErrComponentNotBound))
}

func TestFormatClusterStats(t *testing.T) {
	r := require.New(t)

	t.Run("When not yet loaded and no error, shows a loading state", func(t *testing.T) {
		text := formatClusterStatsFor(&component.ClusterStatsProps{Loaded: false})
		r.Contains(text, "loading")
	})

	t.Run("When not yet loaded and an error occurred, shows the error", func(t *testing.T) {
		text := formatClusterStatsFor(&component.ClusterStatsProps{
			Loaded: false,
			Err:    errors.New("connection refused"),
		})
		r.Contains(text, "connection refused")
	})

	t.Run("When loaded, shows stats", func(t *testing.T) {
		text := formatClusterStatsFor(&component.ClusterStatsProps{
			Loaded: true,
			Stats: &models.ClusterStats{
				Leader: "leader-addr",
				Nodes:  models.NodeStats{Total: 5},
			},
		})
		r.Contains(text, "leader-addr")
		r.Contains(text, "Nodes:   5")
	})

	t.Run("When loaded with a stale error, shows stats plus a stale marker folded into the header (not a separate line)", func(t *testing.T) {
		text := formatClusterStatsFor(&component.ClusterStatsProps{
			Loaded: true,
			Stats:  &models.ClusterStats{Leader: "leader-addr"},
			Err:    errors.New("timeout"),
		})
		r.Contains(text, "leader-addr")
		r.Contains(text, "stale")
	})

	t.Run("When Utilization is set, folds live percentages into the CPU/Memory lines rather than a new line", func(t *testing.T) {
		text := formatClusterStatsFor(&component.ClusterStatsProps{
			Loaded: true,
			Stats: &models.ClusterStats{
				Utilization: &models.ResourceUtilization{CPUPercent: 42, MemoryPercent: 55},
			},
		})
		r.Contains(text, "CPU:")
		r.Contains(text, "(42% live)")
		r.Contains(text, "Memory:")
		r.Contains(text, "(55% live)")
	})

	t.Run("When Utilization is nil, omits the live suffix", func(t *testing.T) {
		text := formatClusterStatsFor(&component.ClusterStatsProps{
			Loaded: true,
			Stats:  &models.ClusterStats{},
		})
		r.NotContains(text, "live)")
	})

	t.Run("Loaded output never exceeds the header's fixed 8-line content budget", func(t *testing.T) {
		// The header row is a fixed 10 lines tall; 8 content lines is what
		// fits. Adding a 9th line (discovered live, against a real
		// cluster) gets silently clipped by tview rather than erroring,
		// so this guards against a regression that would only show up
		// visually.
		text := formatClusterStatsFor(&component.ClusterStatsProps{
			Loaded: true,
			Stats: &models.ClusterStats{
				Leader: "leader-addr",
				Utilization: &models.ResourceUtilization{
					CPUPercent:    42,
					MemoryPercent: 55,
				},
			},
			Err: errors.New("timeout"),
		})

		r.LessOrEqual(len(strings.Split(text, "\n")), 8)
	})
}

// formatClusterStatsFor renders through the exported component so the
// private formatter is exercised without needing to export it.
func formatClusterStatsFor(props *component.ClusterStatsProps) string {
	textView := &componentfakes.FakeTextView{}
	cs := component.NewClusterStats()
	cs.TextView = textView
	cs.Bind(tview.NewFlex())
	cs.Props = props

	cs.Render()

	return textView.SetTextArgsForCall(0)
}
