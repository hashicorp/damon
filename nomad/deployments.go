// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
)

func (n *Nomad) Deployments(so *SearchOptions) ([]*models.Deployment, error) {
	if so == nil {
		so = &SearchOptions{}
	}

	d, _, err := n.DpClient.List(&api.QueryOptions{
		Namespace: so.Namespace,
		Region:    so.Region,
	})

	deps := toDeployments(d)

	return deps, err
}

func toDeployments(dep []*api.Deployment) []*models.Deployment {
	result := make([]*models.Deployment, 0., len(dep))
	for _, d := range dep {
		result = append(result, &models.Deployment{
			ID:                d.ID,
			JobID:             d.JobID,
			Namespace:         d.Namespace,
			Status:            d.Status,
			StatusDescription: d.StatusDescription,
		})

	}
	return result
}
