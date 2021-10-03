package nomad

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
)

func (n *Nomad) JobStatus(jobID string, so *SearchOptions) (*models.JobStatus, error) {
	if so == nil {
		so = &SearchOptions{}
	}

	taskgroups, _ := n.TaskGroups(jobID, so)

	info, _, err := n.JobClient.Info(jobID, &api.QueryOptions{
		Namespace: so.Namespace,
		Region:    so.Region,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve job info: %w", err)
	}

	d, _, err := n.DpClient.List(&api.QueryOptions{
		Namespace: so.Namespace,
		Region:    so.Region,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve job deployments: %w", err)
	}

	taskGroupStatus := toTaskGroupStatus(jobID, d)

	allocations, _ := n.JobAllocs(jobID, so)

	if err != nil {
		return nil, err
	}

	jobStatus := toJobStatus(info, taskgroups, taskGroupStatus, allocations)

	return jobStatus, nil
}

func toJobStatus(job *api.Job, tasks []*models.TaskGroup, taskStatus []*models.TaskGroupStatus, allocs []*models.Alloc) *models.JobStatus {
	jobStatus := &models.JobStatus{
		ID:              *job.ID,
		Name:            *job.Name,
		SubmitDate:      time.Unix(0, *job.SubmitTime),
		Type:            *job.Type,
		Priority:        *job.Priority,
		Datacenters:     strings.Join(job.Datacenters, ", "),
		Namespace:       *job.Namespace,
		Status:          *job.Status,
		Periodic:        job.Periodic != nil,
		Parameterized:   job.ParameterizedJob != nil,
		TaskGroups:      tasks,
		TaskGroupStatus: taskStatus,
		Allocations:     allocs,
	}

	return jobStatus
}

func toTaskGroupStatus(jobID string, dep []*api.Deployment) []*models.TaskGroupStatus {
	result := make([]*models.TaskGroupStatus, 0., len(dep))
	for _, d := range dep {
		if d.JobID == jobID {
			for _, t := range d.TaskGroups {
				result = append(result, &models.TaskGroupStatus{
					ID:                d.JobID,
					Status:            d.Status,
					StatusDescription: d.StatusDescription,
					Desired:           t.DesiredTotal,
					Placed:            t.PlacedAllocs,
					Healthy:           t.HealthyAllocs,
					Unhealthy:         t.UnhealthyAllocs,
					ProgressDeadline:  t.ProgressDeadline,
				})
			}
			break // * Deployments are already sorted. So break once you find your's.
		}

	}
	return result
}
