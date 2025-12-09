// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad

import "github.com/hashicorp/nomad/api"

//TODO fix bug with task name
func (n *Nomad) Logs(allocID, taskName, logType string, cancel <-chan struct{}) (<-chan *api.StreamFrame, <-chan error) {
	return n.AllocFSClient.Logs(
		&api.Allocation{ID: allocID},
		true,
		taskName,
		logType,
		"end",
		20000,
		cancel,
		&api.QueryOptions{},
	)
}
