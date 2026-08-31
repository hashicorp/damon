// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	. "github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
)

func TestClusterStats(t *testing.T) {
	r := require.New(t)

	fakeNodeClient := &nomadfakes.FakeNodeClient{}
	fakeStatusClient := &nomadfakes.FakeStatusClient{}
	fakeAgentClient := &nomadfakes.FakeAgentClient{}

	nomad := &Nomad{
		NodeClient:   fakeNodeClient,
		StatusClient: fakeStatusClient,
		AgentClient:  fakeAgentClient,
	}

	t.Run("When everything is fine", func(t *testing.T) {
		fakeNodeClient.ListReturns([]*api.NodeListStub{
			{ID: "node-1", Status: "ready"},
			{ID: "node-2", Status: "ready"},
		}, nil, nil)
		fakeStatusClient.LeaderReturns("10.0.0.1:4647", nil)
		fakeAgentClient.MembersReturns(&api.ServerMembers{
			Members: []*api.AgentMember{
				{Status: "alive"},
			},
		}, nil)
		fakeNodeClient.InfoReturnsOnCall(0, &api.Node{
			NodeResources: &api.NodeResources{
				Cpu:    api.NodeCpuResources{CpuShares: 2000},
				Memory: api.NodeMemoryResources{MemoryMB: 4096},
			},
		}, nil, nil)
		fakeNodeClient.InfoReturnsOnCall(1, &api.Node{
			NodeResources: &api.NodeResources{
				Cpu:    api.NodeCpuResources{CpuShares: 2000},
				Memory: api.NodeMemoryResources{MemoryMB: 4096},
			},
		}, nil, nil)

		stats, err := nomad.ClusterStats(nil, nil, nil)

		r.NoError(err)
		r.Equal("10.0.0.1:4647", stats.Leader)
		r.Equal(2, stats.Nodes.Total)
		r.Equal(1, stats.Servers.Total)
		r.Equal(4000, stats.Capacity.TotalCPUShares)
		r.Equal(8192, stats.Capacity.TotalMemoryMB)
	})

	t.Run("When a node's Info() call fails, capacity for that node is skipped rather than failing the whole update", func(t *testing.T) {
		fakeNodeClient.ListReturns([]*api.NodeListStub{
			{ID: "node-1", Status: "ready"},
			{ID: "node-2", Status: "ready"},
		}, nil, nil)
		fakeStatusClient.LeaderReturns("10.0.0.1:4647", nil)
		fakeAgentClient.MembersReturns(&api.ServerMembers{}, nil)
		fakeNodeClient.InfoReturnsOnCall(2, nil, nil, errors.New("node unreachable"))
		fakeNodeClient.InfoReturnsOnCall(3, &api.Node{
			NodeResources: &api.NodeResources{
				Cpu:    api.NodeCpuResources{CpuShares: 1000},
				Memory: api.NodeMemoryResources{MemoryMB: 2048},
			},
		}, nil, nil)

		stats, err := nomad.ClusterStats(nil, nil, nil)

		r.NoError(err)
		r.Equal(1000, stats.Capacity.TotalCPUShares)
		r.Equal(2048, stats.Capacity.TotalMemoryMB)
	})

	t.Run("When Info() returns a node with nil NodeResources, it's skipped", func(t *testing.T) {
		fakeNodeClient.ListReturns([]*api.NodeListStub{
			{ID: "node-1", Status: "ready"},
		}, nil, nil)
		fakeStatusClient.LeaderReturns("10.0.0.1:4647", nil)
		fakeAgentClient.MembersReturns(&api.ServerMembers{}, nil)
		fakeNodeClient.InfoReturnsOnCall(4, &api.Node{NodeResources: nil}, nil, nil)

		stats, err := nomad.ClusterStats(nil, nil, nil)

		r.NoError(err)
		r.Equal(0, stats.Capacity.TotalCPUShares)
	})

	t.Run("When Nodes().List() fails", func(t *testing.T) {
		fakeNodeClient.ListReturns(nil, nil, errors.New("nodes failed"))

		stats, err := nomad.ClusterStats(nil, nil, nil)

		r.Nil(stats)
		r.EqualError(err, "nodes failed")
	})

	t.Run("When Status().Leader() fails", func(t *testing.T) {
		fakeNodeClient.ListReturns([]*api.NodeListStub{}, nil, nil)
		fakeStatusClient.LeaderReturns("", errors.New("leader failed"))

		stats, err := nomad.ClusterStats(nil, nil, nil)

		r.Nil(stats)
		r.EqualError(err, "leader failed")
	})

	t.Run("When Agent().Members() fails", func(t *testing.T) {
		fakeStatusClient.LeaderReturns("10.0.0.1:4647", nil)
		fakeAgentClient.MembersReturns(nil, errors.New("members failed"))

		stats, err := nomad.ClusterStats(nil, nil, nil)

		r.Nil(stats)
		r.EqualError(err, "members failed")
	})
}
