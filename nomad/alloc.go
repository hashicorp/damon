package nomad

import (
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

	allocs := toAllocs(list)

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

	allocs := toAllocs(list)

	return allocs, nil
}

func toAllocs(list []*api.AllocationListStub) []*models.Alloc {
	result := make([]*models.Alloc, 0, len(list))
	for _, el := range list {
		alloc := &models.Alloc{
			ID:            el.ID,
			TaskGroup:     el.TaskGroup,
			JobID:         el.JobID,
			JobType:       el.JobType,
			NodeID:        el.NodeID,
			NodeName:      el.NodeName,
			DesiredStatus: el.DesiredStatus,
		}

		for k := range el.TaskStates {
			alloc.TaskNames = append(alloc.TaskNames, k)
		}

		result = append(result, alloc)
	}

	return result
}
