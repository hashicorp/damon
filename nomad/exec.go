// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad

import (
	"context"
	"io"

	"github.com/hashicorp/nomad/api"
)

// Exec runs command inside the given task of the given allocation, attaching
// stdin/stdout/stderr to it as a TTY session. It blocks until the command
// terminates and returns its exit code.
func (n *Nomad) Exec(ctx context.Context, allocID, taskName string, command []string,
	stdin io.Reader, stdout, stderr io.Writer, resizeCh <-chan api.TerminalSize) (int, error) {
	return n.AllocClient.Exec(
		ctx,
		&api.Allocation{ID: allocID},
		taskName,
		true,
		command,
		stdin,
		stdout,
		stderr,
		resizeCh,
		&api.QueryOptions{},
	)
}
