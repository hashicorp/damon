// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package models_test

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/models"
)

func TestNewClusterStats(t *testing.T) {
	r := require.New(t)

	nodes := []*api.NodeListStub{
		{
			Status: "ready",
			NodeResources: &api.NodeResources{
				Cpu:    api.NodeCpuResources{CpuShares: 2000},
				Memory: api.NodeMemoryResources{MemoryMB: 4096},
			},
		},
		{
			Status: "ready",
			NodeResources: &api.NodeResources{
				Cpu:    api.NodeCpuResources{CpuShares: 2000},
				Memory: api.NodeMemoryResources{MemoryMB: 4096},
			},
		},
		{
			Status: "down",
		},
	}

	members := []*api.AgentMember{
		{Status: "alive"},
		{Status: "alive"},
		{Status: "failed"},
	}

	jobs := []*models.Job{
		{Status: "running"},
		{Status: "running"},
		{Status: "dead"},
	}

	allocs := []*models.Alloc{
		{Status: "running", AllocatedCPU: 500, AllocatedMemMB: 256},
		{Status: "running", AllocatedCPU: 300, AllocatedMemMB: 128},
		{Status: "failed", AllocatedCPU: 100, AllocatedMemMB: 64},
	}

	t.Run("aggregates node, server, job, and alloc counts", func(t *testing.T) {
		stats := models.NewClusterStats(nodes, "10.0.0.1:4647", members, jobs, allocs)

		r.Equal("10.0.0.1:4647", stats.Leader)

		r.Equal(3, stats.Nodes.Total)
		r.Equal(2, stats.Nodes.ByStatus["ready"])
		r.Equal(1, stats.Nodes.ByStatus["down"])

		r.Equal(3, stats.Servers.Total)
		r.Equal(2, stats.Servers.ByStatus["alive"])
		r.Equal(1, stats.Servers.ByStatus["failed"])

		r.Equal(3, stats.Jobs.Total)
		r.Equal(2, stats.Jobs.ByStatus["running"])
		r.Equal(1, stats.Jobs.ByStatus["dead"])

		r.Equal(3, stats.Allocations.Total)
		r.Equal(2, stats.Allocations.ByStatus["running"])
		r.Equal(1, stats.Allocations.ByStatus["failed"])
	})

	t.Run("does not compute total capacity from the node list (unreliable on real Nomad servers)", func(t *testing.T) {
		stats := models.NewClusterStats(nodes, "", nil, nil, nil)

		r.Equal(0, stats.Capacity.TotalCPUShares)
		r.Equal(0, stats.Capacity.TotalMemoryMB)
	})

	t.Run("sums allocated CPU/memory across allocations", func(t *testing.T) {
		stats := models.NewClusterStats(nil, "", nil, nil, allocs)

		r.Equal(900, stats.Capacity.AllocatedCPU)
		r.Equal(448, stats.Capacity.AllocatedMemoryMB)
	})

	t.Run("handles empty input without panicking", func(t *testing.T) {
		stats := models.NewClusterStats(nil, "", nil, nil, nil)

		r.Equal(0, stats.Nodes.Total)
		r.Equal(0, stats.Servers.Total)
		r.Equal(0, stats.Jobs.Total)
		r.Equal(0, stats.Allocations.Total)
		r.Nil(stats.Utilization)
	})
}

func TestFindHostPort(t *testing.T) {
	r := require.New(t)

	allocs := []*models.Alloc{
		{
			JobID:  "prometheus",
			Status: models.StatusRunning,
			HostPorts: []models.HostPort{
				{Label: "metrics", HostIP: "10.0.0.5", Port: 8080},
			},
		},
		{
			JobID:  "prometheus",
			Status: models.StatusRunning,
			HostPorts: []models.HostPort{
				{Label: "http", HostIP: "10.0.0.9", Port: 9090},
			},
		},
		{
			JobID:  "prometheus",
			Status: models.StatusDead,
			HostPorts: []models.HostPort{
				{Label: "http", HostIP: "10.0.0.1", Port: 9091},
			},
		},
		{
			JobID:  "other-job",
			Status: models.StatusRunning,
			HostPorts: []models.HostPort{
				{Label: "http", HostIP: "10.0.0.2", Port: 9092},
			},
		},
	}

	t.Run("finds the running allocation of the job with the matching port label", func(t *testing.T) {
		addr, ok := models.FindHostPort(allocs, "prometheus", "http")

		r.True(ok)
		r.Equal("10.0.0.9:9090", addr)
	})

	t.Run("ignores allocations that aren't running", func(t *testing.T) {
		_, ok := models.FindHostPort([]*models.Alloc{allocs[2]}, "prometheus", "http")

		r.False(ok)
	})

	t.Run("ignores allocations of other jobs", func(t *testing.T) {
		_, ok := models.FindHostPort([]*models.Alloc{allocs[3]}, "prometheus", "http")

		r.False(ok)
	})

	t.Run("returns false when no allocation exposes the port label", func(t *testing.T) {
		_, ok := models.FindHostPort(allocs, "prometheus", "grpc")

		r.False(ok)
	})

	t.Run("returns false for empty input", func(t *testing.T) {
		_, ok := models.FindHostPort(nil, "prometheus", "http")

		r.False(ok)
	})
}
