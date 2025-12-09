// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"context"

	"github.com/hashicorp/nomad/api"
)

type Topics map[api.Topic][]string

func (n *Nomad) Stream(topics Topics, index uint64) (<-chan *api.Events, error) {
	return n.EventsClient.Stream(context.Background(), topics, index, nil)
}
