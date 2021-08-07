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
		job := toJob(j)
		jobs = append(jobs, job)
	}

	return jobs, err
}

func toJob(j *api.JobListStub) *models.Job {
	t := time.Unix(0, j.SubmitTime)

	return &models.Job{
		ID:                j.ID,
		Name:              j.Name,
		Namespace:         j.JobSummary.Namespace,
		Type:              j.Type,
		Status:            j.Status,
		StatusDescription: j.StatusDescription,
		SubmitTime:        t,
	}
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
