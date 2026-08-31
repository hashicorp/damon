// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package watcher_test

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/state"
	"github.com/hcjulz/damon/watcher"
	"github.com/hcjulz/damon/watcher/watcherfakes"
)

func TestWatchClusterStats_Happy(t *testing.T) {
	r := require.New(t)

	t.Run("It notifies the subscriber initially", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		s := state.New()
		w := watcher.NewWatcher(s, nomad, time.Millisecond*250)

		nomad.ClusterStatsReturnsOnCall(0, &models.ClusterStats{Leader: "first"}, nil)
		nomad.ClusterStatsReturnsOnCall(1, &models.ClusterStats{Leader: "second"}, nil)

		done := make(chan struct{})
		var callCount int
		notifier := func() {
			callCount++
			switch callCount {
			case 1:
				r.Equal("first", s.ClusterStats.Leader)
			case 2:
				done <- struct{}{}
			}
		}

		w.WatchClusterStats(time.Millisecond*50, notifier)

		<-done
	})

	t.Run("It continues to notify the subscriber after the initial notification", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		s := state.New()
		w := watcher.NewWatcher(s, nomad, time.Millisecond*250)

		nomad.ClusterStatsReturnsOnCall(2, &models.ClusterStats{Leader: "third"}, nil)

		done := make(chan struct{})
		var callCount int
		notifier := func() {
			callCount++
			if callCount == 3 {
				defer func() { done <- struct{}{} }()
				r.Equal("third", s.ClusterStats.Leader)
			}
		}

		w.WatchClusterStats(time.Millisecond*50, notifier)

		<-done
		r.Equal(3, callCount)
	})
}

func TestWatchClusterStats_Sad(t *testing.T) {
	r := require.New(t)

	nomad := &watcherfakes.FakeNomad{}
	s := state.New()
	w := watcher.NewWatcher(s, nomad, time.Millisecond*250)

	var modalCalled bool
	w.SubscribeHandler(models.HandleError, func(_ string, _ ...interface{}) {
		modalCalled = true
	})

	nomad.ClusterStatsReturns(nil, errors.New("argh"))

	w.WatchClusterStats(time.Millisecond*250, func() {})

	r.EqualError(s.ClusterStatsErr, "argh")
	r.False(modalCalled, "cluster stats errors must not trigger the blocking error modal")
}

func TestWatchClusterStats_WithPrometheus(t *testing.T) {
	r := require.New(t)

	t.Run("When Prometheus is configured and succeeds, utilization is attached", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		prom := &watcherfakes.FakePrometheusClient{}
		s := state.New()
		w := watcher.NewWatcher(s, nomad, time.Millisecond*250)
		w.SetPrometheusClient(prom)

		nomad.ClusterStatsReturns(&models.ClusterStats{Leader: "leader"}, nil)
		prom.ClusterUtilizationReturns(&models.ResourceUtilization{CPUPercent: 42, MemoryPercent: 55}, nil)

		w.WatchClusterStats(time.Millisecond*250, func() {})

		r.NotNil(s.ClusterStats.Utilization)
		r.Equal(42.0, s.ClusterStats.Utilization.CPUPercent)
	})

	t.Run("When Prometheus fails, base cluster stats still populate", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		prom := &watcherfakes.FakePrometheusClient{}
		s := state.New()
		w := watcher.NewWatcher(s, nomad, time.Millisecond*250)
		w.SetPrometheusClient(prom)

		nomad.ClusterStatsReturns(&models.ClusterStats{Leader: "leader"}, nil)
		prom.ClusterUtilizationReturns(nil, errors.New("prometheus unreachable"))

		w.WatchClusterStats(time.Millisecond*250, func() {})

		r.Equal("leader", s.ClusterStats.Leader)
		r.Nil(s.ClusterStats.Utilization)
	})
}

func TestWatchClusterStats_WithPrometheusDiscovery(t *testing.T) {
	r := require.New(t)

	t.Run("When a running allocation of the job exposes the port label, utilization is queried from it", func(t *testing.T) {
		promServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			query := req.URL.Query().Get("query")
			switch query {
			case "avg(nomad_client_host_cpu_total_percent)":
				w.Write([]byte(`{"status":"success","data":{"result":[{"value":[1234,"33.0"]}]}}`))
			case "100 * avg(nomad_client_host_memory_used) / avg(nomad_client_host_memory_total)":
				w.Write([]byte(`{"status":"success","data":{"result":[{"value":[1234,"44.0"]}]}}`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer promServer.Close()

		promURL, err := url.Parse(promServer.URL)
		r.NoError(err)
		host, portStr, err := net.SplitHostPort(promURL.Host)
		r.NoError(err)
		port, err := strconv.Atoi(portStr)
		r.NoError(err)

		nomad := &watcherfakes.FakeNomad{}
		s := state.New()
		s.Allocations = []*models.Alloc{
			{
				JobID:  "prometheus",
				Status: models.StatusRunning,
				HostPorts: []models.HostPort{
					{Label: "http", HostIP: host, Port: port},
				},
			},
		}
		w := watcher.NewWatcher(s, nomad, time.Millisecond*250)
		w.SetPrometheusDiscovery("prometheus", "http")

		nomad.ClusterStatsReturns(&models.ClusterStats{Leader: "leader"}, nil)

		w.WatchClusterStats(time.Millisecond*250, func() {})

		r.NotNil(s.ClusterStats.Utilization)
		r.Equal(33.0, s.ClusterStats.Utilization.CPUPercent)
		r.Equal(44.0, s.ClusterStats.Utilization.MemoryPercent)
	})

	t.Run("When no running allocation of the job exposes the port label, base cluster stats still populate", func(t *testing.T) {
		nomad := &watcherfakes.FakeNomad{}
		s := state.New()
		w := watcher.NewWatcher(s, nomad, time.Millisecond*250)
		w.SetPrometheusDiscovery("prometheus", "http")

		nomad.ClusterStatsReturns(&models.ClusterStats{Leader: "leader"}, nil)

		w.WatchClusterStats(time.Millisecond*250, func() {})

		r.Equal("leader", s.ClusterStats.Leader)
		r.Nil(s.ClusterStats.Utilization)
	})
}
