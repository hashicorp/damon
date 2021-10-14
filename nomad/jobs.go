package nomad

import (
	"fmt"
	"time"

	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
)

func (n *Nomad) Jobs(so *SearchOptions) ([]*models.Job, error) {
	if so == nil {
		so = &SearchOptions{}
	}

	jobList, _, err := n.JobClient.List(&api.QueryOptions{
		Namespace: so.Namespace,
		Region:    so.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve job list: %w", err)
	}

	var jobs []*models.Job
	for _, j := range jobList {
		readyStatus, deploymentStatus := n.jobSummary(j.ID)
		job := toJob(j)
		job.ReadyStatus = readyStatus
		job.DeploymentStatus = deploymentStatus
		jobs = append(jobs, job)
	}

	return jobs, err
}

func toJob(j *api.JobListStub) *models.Job {
	t := time.Unix(0, j.SubmitTime)

	total := len(j.JobSummary.Summary)
	summary := models.Summary{
		Total: total,
	}

	for _, job := range j.JobSummary.Summary {
		if job.Running > 0 {
			summary.Running++
		}
	}

	return &models.Job{
		ID:                j.ID,
		Name:              j.Name,
		Namespace:         j.JobSummary.Namespace,
		Type:              j.Type,
		Status:            j.Status,
		StatusDescription: j.StatusDescription,
		StatusSummary:     summary,
		SubmitTime:        t,
	}
}

func (n *Nomad) jobSummary(jobID string) (models.ReadyStatus, string) {
	taskGroupStatus, deploymentStatus, _ := n.TaskStatus(jobID)
	status := models.ReadyStatus{
		Desired:   0,
		Running:   0,
		Healthy:   0,
		Unhealthy: 0,
	}
	for _, taskGroup := range taskGroupStatus {
		status.Desired += taskGroup.Desired
		status.Running += taskGroup.Placed
		status.Healthy += taskGroup.Healthy
		status.Unhealthy += taskGroup.Unhealthy
	}
	return status, deploymentStatus
}

func (n *Nomad) GetJob(jobID string) (*api.Job, error) {
	job, _, err := n.JobClient.Info(jobID, nil)
	return job, err
}

func (n *Nomad) StartJob(job *api.Job) error {
	stop := false
	job.Stop = &stop

	_, _, err := n.JobClient.Register(job, nil)
	return err
}

func (n *Nomad) StopJob(jobID string) error {
	_, _, err := n.JobClient.Deregister(jobID, false, nil)
	return err
}
