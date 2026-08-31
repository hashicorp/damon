// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package models

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type ClusterStats struct {
	Leader      string
	Nodes       NodeStats
	Servers     ServerStats
	Jobs        JobCounts
	Allocations AllocCounts
	Capacity    ResourceStats
	Utilization *ResourceUtilization
}

type NodeStats struct {
	Total    int
	ByStatus map[string]int
}

type ServerStats struct {
	Total    int
	ByStatus map[string]int
}

type JobCounts struct {
	Total    int
	ByStatus map[string]int
}

type AllocCounts struct {
	Total    int
	ByStatus map[string]int
}

type ResourceStats struct {
	TotalCPUShares    int
	TotalMemoryMB     int
	AllocatedCPU      int
	AllocatedMemoryMB int
}

// ResourceUtilization holds live cluster-average utilization sourced
// from an external Prometheus server. It's only populated when a
// Prometheus integration is configured, so it's optional everywhere
// it appears.
type ResourceUtilization struct {
	CPUPercent    float64
	MemoryPercent float64
}

// NewClusterStats aggregates cluster-wide stats from data already
// fetched elsewhere (node list, leader, server members, jobs, and
// allocations). It is a pure function with no Nomad-SDK network calls
// of its own, so the counting/summing logic can be unit tested with
// plain fixtures.
//
// Total node capacity (Capacity.TotalCPUShares/TotalMemoryMB) is
// deliberately NOT computed here from nodes[i].NodeResources: on real
// Nomad servers that field is only populated by the per-node info
// endpoint, not the list endpoint used to build `nodes`. Callers that
// want real totals fetch per-node info separately and set
// stats.Capacity.Total* after calling this function.
func NewClusterStats(
	nodes []*api.NodeListStub,
	leader string,
	members []*api.AgentMember,
	jobs []*Job,
	allocs []*Alloc,
) *ClusterStats {
	nodeStats := NodeStats{ByStatus: map[string]int{}}
	capacity := ResourceStats{}

	for _, n := range nodes {
		nodeStats.Total++
		nodeStats.ByStatus[n.Status]++
	}

	serverStats := ServerStats{ByStatus: map[string]int{}}
	for _, m := range members {
		serverStats.Total++
		serverStats.ByStatus[m.Status]++
	}

	jobCounts := JobCounts{ByStatus: map[string]int{}}
	for _, j := range jobs {
		jobCounts.Total++
		jobCounts.ByStatus[j.Status]++
	}

	allocCounts := AllocCounts{ByStatus: map[string]int{}}
	for _, a := range allocs {
		allocCounts.Total++
		allocCounts.ByStatus[a.Status]++
		capacity.AllocatedCPU += a.AllocatedCPU
		capacity.AllocatedMemoryMB += a.AllocatedMemMB
	}

	return &ClusterStats{
		Leader:      leader,
		Nodes:       nodeStats,
		Servers:     serverStats,
		Jobs:        jobCounts,
		Allocations: allocCounts,
		Capacity:    capacity,
	}
}

// FindHostPort looks for a running allocation of jobID exposing
// portLabel, and returns "host:port" for it. It searches allocs
// directly (typically state.Allocations, which the app already keeps
// fresh across all namespaces) rather than making a fresh Nomad call,
// so this is free to call on every poll. ok is false if no match is
// found.
func FindHostPort(allocs []*Alloc, jobID, portLabel string) (addr string, ok bool) {
	for _, a := range allocs {
		if a.JobID != jobID || a.Status != StatusRunning {
			continue
		}

		for _, p := range a.HostPorts {
			if p.Label == portLabel {
				return fmt.Sprintf("%s:%d", p.HostIP, p.Port), true
			}
		}
	}

	return "", false
}
