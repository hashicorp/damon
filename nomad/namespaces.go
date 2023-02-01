// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"github.com/hashicorp/nomad/api"

	"github.com/hcjulz/damon/models"
)

func (n *Nomad) Namespaces(_ *SearchOptions) ([]*models.Namespace, error) {
	ns, _, err := n.NsClient.List(nil)
	if err != nil {
		return nil, err
	}

	namespaces := []*models.Namespace{}
	for _, s := range ns {
		namespaces = append(namespaces, toNamespace(s))
	}

	return namespaces, nil
}

func toNamespace(j *api.Namespace) *models.Namespace {
	return &models.Namespace{
		Name:        j.Name,
		Description: j.Description,
	}
}
