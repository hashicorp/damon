// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	primitive "github.com/hcjulz/damon/primitives"

	"github.com/hcjulz/damon/models"
)

type ClusterStats struct {
	TextView TextView
	Props    *ClusterStatsProps

	slot *tview.Flex
}

type ClusterStatsProps struct {
	Stats  *models.ClusterStats
	Err    error
	Loaded bool
}

func NewClusterStats() *ClusterStats {
	return &ClusterStats{
		TextView: primitive.NewTextView(tview.AlignRight),
		Props:    &ClusterStatsProps{},
	}
}

func (c *ClusterStats) Render() error {
	if c.slot == nil {
		return ErrComponentNotBound
	}

	c.TextView.SetText(formatClusterStats(c.Props))

	c.slot.Clear()
	c.slot.AddItem(c.TextView.Primitive(), 0, 1, false)

	return nil
}

func (c *ClusterStats) Bind(slot *tview.Flex) {
	c.slot = slot
}

func formatClusterStats(props *ClusterStatsProps) string {
	if !props.Loaded {
		if props.Err != nil {
			return fmt.Sprintf("[red]cluster stats unavailable:\n%s", props.Err)
		}
		return "[gray]loading cluster stats..."
	}

	stats := props.Stats

	// The header row is a fixed height, so lines are kept to a strict
	// budget rather than growing with optional data — live utilization
	// (when available) is folded into the existing CPU/Memory lines
	// instead of adding a 9th line, which tview silently clips.
	var cpuLive, memLive string
	if stats.Utilization != nil {
		cpuLive = fmt.Sprintf(" (%.0f%% live)", stats.Utilization.CPUPercent)
		memLive = fmt.Sprintf(" (%.0f%% live)", stats.Utilization.MemoryPercent)
	}

	// A stale-data indicator is folded into the header line rather than
	// appended as its own line, for the same fixed-height reason.
	header := "[#26ffe6]Cluster"
	if props.Err != nil {
		header = "[#26ffe6]Cluster [red](stale)"
	}

	lines := []string{
		header,
		fmt.Sprintf("Leader:  %s", stats.Leader),
		fmt.Sprintf("Nodes:   %d", stats.Nodes.Total),
		fmt.Sprintf("Servers: %d", stats.Servers.Total),
		fmt.Sprintf("Jobs:    %d", stats.Jobs.Total),
		fmt.Sprintf("Allocs:  %d", stats.Allocations.Total),
		fmt.Sprintf("CPU:     %d/%s MHz%s", stats.Capacity.AllocatedCPU, capacityOrNA(stats.Capacity.TotalCPUShares), cpuLive),
		fmt.Sprintf("Memory:  %d/%s MB%s", stats.Capacity.AllocatedMemoryMB, capacityOrNA(stats.Capacity.TotalMemoryMB), memLive),
	}

	return strings.Join(lines, "\n")
}

// capacityOrNA renders "n/a" for zero total capacity. Some Nomad
// servers don't populate NodeResources on the node list endpoint
// (only on the per-node info endpoint), so a total of 0 means
// "unknown," not "no capacity."
func capacityOrNA(total int) string {
	if total == 0 {
		return "n/a"
	}
	return fmt.Sprint(total)
}
