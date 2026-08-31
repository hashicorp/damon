// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
)

// ClusterStats aggregates cheap, cluster-wide Nomad data: node counts,
// the current leader, server/member counts, and job/alloc resource
// aggregates derived from data the caller already has (allocs, jobs).
// It deliberately takes allocs/jobs as parameters rather than fetching
// them itself, since they're already kept fresh elsewhere (the
// long-lived event-stream watcher) and re-fetching them here would
// duplicate that cost.
func (n *Nomad) ClusterStats(allocs []*models.Alloc, jobs []*models.Job, so *SearchOptions) (*models.ClusterStats, error) {
	if so == nil {
		so = &SearchOptions{}
	}

	nodes, _, err := n.NodeClient.List(nil)
	if err != nil {
		return nil, err
	}

	leader, err := n.StatusClient.Leader()
	if err != nil {
		return nil, err
	}

	serverMembers, err := n.AgentClient.Members()
	if err != nil {
		return nil, err
	}

	var members []*api.AgentMember
	if serverMembers != nil {
		members = serverMembers.Members
	}

	// The node list endpoint doesn't populate NodeResources on every
	// Nomad version (confirmed empty on a live 1.x/2.x server despite
	// the SDK type declaring the field) — only the per-node info
	// endpoint reliably does. Capacity is otherwise-static data
	// (changes only when nodes join/leave or hardware changes), so
	// paying one Info() call per node here is an accepted tradeoff to
	// get real numbers instead of "n/a".
	totalCPU, totalMem := n.nodeCapacity(nodes)

	stats := models.NewClusterStats(nodes, leader, members, jobs, allocs)
	stats.Capacity.TotalCPUShares = totalCPU
	stats.Capacity.TotalMemoryMB = totalMem

	return stats, nil
}

func (n *Nomad) nodeCapacity(nodes []*api.NodeListStub) (totalCPU, totalMem int) {
	for _, node := range nodes {
		info, _, err := n.NodeClient.Info(node.ID, nil)
		if err != nil || info.NodeResources == nil {
			continue
		}

		totalCPU += int(info.NodeResources.Cpu.CpuShares)
		totalMem += int(info.NodeResources.Memory.MemoryMB)
	}

	return totalCPU, totalMem
}
