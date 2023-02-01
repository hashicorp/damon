// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"fmt"
	"time"

	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
)

func (n *Nomad) JobAllocs(jobID string, so *SearchOptions) ([]*models.Alloc, error) {
	if so == nil {
		so = &SearchOptions{}
	}

	list, _, err := n.JobClient.Allocations(jobID, false, &api.QueryOptions{
		Namespace: so.Namespace,
		Region:    so.Region,
	})

	if err != nil {
		return nil, err
	}

	allocs, err := n.toAllocs(list)
	if err != nil {
		return nil, err
	}

	return allocs, nil
}

func (n *Nomad) Allocations(so *SearchOptions) ([]*models.Alloc, error) {
	if so == nil {
		so = &SearchOptions{}
	}

	list, _, err := n.AllocClient.List(&api.QueryOptions{
		Namespace: so.Namespace,
		Region:    so.Region,
	})
	if err != nil {
		return nil, err
	}

	allocs, err := n.toAllocs(list)
	if err != nil {
		return nil, err
	}

	return allocs, nil
}

func getTasksFromAlloc(taskStates map[string]*api.TaskState, alloc *api.Allocation) []*models.Task {
	tasks := []*models.Task{}

	for _, t := range alloc.GetTaskGroup().Tasks {
		task := &models.Task{
			Name:     t.Name,
			State:    taskStates[t.Name].State,
			Events:   taskStates[t.Name].Events,
			Driver:   t.Driver,
			Env:      t.Env,
			Config:   t.Config,
			CPU:      *t.Resources.CPU,
			MemoryMB: *t.Resources.MemoryMB,
		}

		tasks = append(tasks, task)

	}

	return tasks
}

func (n *Nomad) toAllocs(list []*api.AllocationListStub) ([]*models.Alloc, error) {
	result := make([]*models.Alloc, 0, len(list))
	for _, el := range list {
		a, _, err := n.AllocClient.Info(el.ID, &api.QueryOptions{
			Namespace: "*",
		})
		if err != nil {
			return nil, err
		}

		tasks := getTasksFromAlloc(el.TaskStates, a)

		alloc := &models.Alloc{
			ID:            el.ID,
			Namespace:     el.Namespace,
			TaskGroup:     el.TaskGroup,
			TaskList:      tasks,
			JobID:         el.JobID,
			JobType:       el.JobType,
			NodeID:        el.NodeID,
			NodeName:      el.NodeName,
			DesiredStatus: el.DesiredStatus,
			Version:       el.JobVersion,
			Status:        el.ClientStatus,
			Created:       time.Unix(0, el.CreateTime),
			Modified:      time.Unix(0, el.ModifyTime),
		}

		for k, t := range el.TaskStates {
			alloc.TaskNames = append(alloc.TaskNames, k)
			alloc.Tasks = append(alloc.Tasks, models.AllocTask{
				Name:   k,
				Events: t.Events,
			})
		}

		if a.AllocatedResources != nil {
			for _, net := range a.AllocatedResources.Shared.Ports {
				alloc.HostAddresses = append(
					alloc.HostAddresses,
					fmt.Sprintf("%s/%s:%d", net.Label, net.HostIP, net.Value),
				)
			}

		}

		result = append(result, alloc)
	}

	return result, nil
}
