package nomad

import (
	"sort"

	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
)

func (n *Nomad) TaskGroups(jobID string, so *SearchOptions) ([]*models.TaskGroup, error) {
	if so == nil {
		so = &SearchOptions{}
	}

	summary, _, err := n.JobClient.Summary(jobID, &api.QueryOptions{
		Namespace: so.Namespace,
		Region:    so.Region,
	})

	taskGroups := toTaskGroups(summary)

	return taskGroups, err
}

func toTaskGroups(js *api.JobSummary) []*models.TaskGroup {
	if js == nil {
		return nil
	}

	keys := sortTaskGroupSummaryByKey(js.Summary)

	var result []*models.TaskGroup
	for _, key := range keys {
		tgs := js.Summary[key]
		taskGroup := &models.TaskGroup{
			Name:     key,
			JobID:    js.JobID,
			Queued:   tgs.Queued,
			Complete: tgs.Complete,
			Failed:   tgs.Failed,
			Running:  tgs.Running,
			Starting: tgs.Starting,
			Lost:     tgs.Lost,
		}

		result = append(result, taskGroup)
	}

	return result
}

func sortTaskGroupSummaryByKey(tgs map[string]api.TaskGroupSummary) []string {
	keys := make([]string, 0, len(tgs))
	for k := range tgs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}
