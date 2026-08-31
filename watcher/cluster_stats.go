// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package watcher

import (
	"fmt"
	"time"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/prometheus"
)

// WatchClusterStats starts an independent polling loop for
// cluster-wide stats (nodes, leader, members, job/alloc/resource
// aggregates). Unlike SubscribeToX methods, this is not torn down by
// Subscribe/Unsubscribe/activities.DeactivateAll: it runs for the
// entire application lifetime, because the stats panel lives in the
// always-visible header, not the swappable body.
func (w *Watcher) WatchClusterStats(interval time.Duration, notify func()) {
	w.updateClusterStats()
	notify()

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			w.updateClusterStats()
			notify()
		}
	}()
}

func (w *Watcher) updateClusterStats() {
	stats, err := w.nomad.ClusterStats(w.state.Allocations, w.state.Jobs, nil)
	if err != nil || stats == nil {
		// Deliberately not routed through NotifyHandler(HandleError, ...):
		// that pops the blocking error modal, which would interrupt
		// whatever view the user is on every time a background header
		// poll fails. Keep the last-good stats in place and surface the
		// error quietly in the panel instead.
		if err == nil {
			err = fmt.Errorf("nomad returned no cluster stats")
		}
		w.state.ClusterStatsErr = err
		return
	}

	w.state.ClusterStatsErr = nil
	stats.Utilization = w.resolveUtilization()
	w.state.ClusterStats = stats
}

// resolveUtilization returns live utilization from whichever
// Prometheus source is configured (a static client, or a Nomad-based
// discovery re-resolved on every call so it self-heals if the
// Prometheus job moves), or nil if neither is configured or the
// lookup/query fails. Failures here never fail the broader cluster
// stats update — utilization is a bonus metric, not a core stat.
func (w *Watcher) resolveUtilization() *models.ResourceUtilization {
	client := w.prometheus

	if client == nil && w.promDisc != nil {
		addr, ok := models.FindHostPort(w.state.Allocations, w.promDisc.jobID, w.promDisc.portLabel)
		if !ok {
			return nil
		}
		client = prometheus.New("http://" + addr)
	}

	if client == nil {
		return nil
	}

	util, err := client.ClusterUtilization()
	if err != nil {
		return nil
	}

	return util
}
