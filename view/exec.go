// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/nomad/api"
	"golang.org/x/term"
)

// Exec suspends the TUI and opens an interactive shell inside the given
// task/allocation, equivalent to `nomad alloc exec -task <task> <allocID> /bin/sh`.
func (v *View) Exec(taskName, allocID string) {
	var execErr error

	// The background watcher keeps polling/streaming while the screen is
	// suspended and would otherwise call Draw() on every refresh - drawing
	// to a screen that's mid-handoff to the exec session's raw terminal is
	// unsafe and can leave tview's internal state out of sync with the
	// actual terminal once we resume, which is what caused the TUI to come
	// back blank and unresponsive after exiting the remote shell.
	v.Watcher.Unsubscribe()

	v.Layout.Container.Suspend(func() {
		stdinFd := int(os.Stdin.Fd())

		// Suspend() only restores the terminal to its normal ("cooked") mode,
		// the same as dropping to an ordinary shell. In cooked mode the local
		// tty driver intercepts control characters (Ctrl-C, Ctrl-D, ...) and
		// turns them into local signals instead of forwarding them as literal
		// bytes for the remote shell to interpret - e.g. Ctrl-C would kill
		// Damon itself rather than interrupting the remote command. Put the
		// local terminal in raw mode for the duration of the session, like
		// `nomad alloc exec`/`docker exec -it` do.
		if oldState, err := term.MakeRaw(stdinFd); err == nil {
			defer term.Restore(stdinFd, oldState)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		resizeCh := make(chan api.TerminalSize, 1)
		if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			resizeCh <- api.TerminalSize{Width: w, Height: h}
		}

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGWINCH)
		defer signal.Stop(sigCh)

		go func() {
			for range sigCh {
				if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
					select {
					case resizeCh <- api.TerminalSize{Width: w, Height: h}:
					default:
					}
				}
			}
		}()

		_, execErr = v.Client.Exec(ctx, allocID, taskName, []string{"/bin/sh"},
			stripCPR(os.Stdin), os.Stdout, os.Stderr, resizeCh)
	})

	if execErr != nil {
		v.handleError("exec failed: %s", execErr.Error())
	}

	// Force a full re-sync in addition to the normal draw: after a suspend,
	// the screen buffer was torn down and redrawn from scratch, so a plain
	// incremental Draw() isn't guaranteed to leave the terminal in sync with
	// what tview thinks is on it.
	v.Layout.Container.Sync()

	// Re-render the Tasks view we came from: this both repaints known-good
	// content (rather than relying on whatever was left over from before
	// the suspend) and re-subscribes to the watcher we unsubscribed from
	// above.
	if alloc, ok := v.getAllocation(allocID); ok {
		v.Tasks(alloc)
		return
	}

	v.Draw()
}
