// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/config"
)

func writeFile(t *testing.T, dir, contents string) string {
	t.Helper()
	path := filepath.Join(dir, "damon.yaml")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o644))
	return path
}

func TestLoadFrom(t *testing.T) {
	r := require.New(t)

	t.Run("When no file exists at either path, returns defaults", func(t *testing.T) {
		dir := t.TempDir()
		cfg, err := config.LoadFrom(filepath.Join(dir, "missing.yaml"), filepath.Join(dir, "also-missing.yaml"))

		r.NoError(err)
		r.Equal(config.DefaultClusterStatsRefreshInterval, cfg.ClusterStatsRefreshInterval)
		r.Empty(cfg.PrometheusURL)
		r.Equal(config.PrometheusPlacementNone, cfg.PrometheusPlacement)
		r.Equal(config.DefaultPrometheusJob, cfg.PrometheusJob)
		r.Equal(config.DefaultPrometheusPortLabel, cfg.PrometheusPortLabel)
	})

	t.Run("When only the cwd path exists, uses it", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "cluster_stats_refresh_interval: 5s\n")

		cfg, err := config.LoadFrom(path, filepath.Join(dir, "missing.yaml"))

		r.NoError(err)
		r.Equal(5*time.Second, cfg.ClusterStatsRefreshInterval)
	})

	t.Run("When only the home path exists, uses it", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "cluster_stats_refresh_interval: 15s\n")

		cfg, err := config.LoadFrom(filepath.Join(dir, "missing.yaml"), path)

		r.NoError(err)
		r.Equal(15*time.Second, cfg.ClusterStatsRefreshInterval)
	})

	t.Run("When both paths exist, cwd wins", func(t *testing.T) {
		cwdDir := t.TempDir()
		homeDir := t.TempDir()

		cwdPath := writeFile(t, cwdDir, "cluster_stats_refresh_interval: 5s\n")
		homePath := writeFile(t, homeDir, "cluster_stats_refresh_interval: 15s\n")

		cfg, err := config.LoadFrom(cwdPath, homePath)

		r.NoError(err)
		r.Equal(5*time.Second, cfg.ClusterStatsRefreshInterval)
	})

	t.Run("When the file has malformed YAML, returns an error", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "not: valid: yaml: [")

		_, err := config.LoadFrom(path, "")

		r.Error(err)
	})

	t.Run("When the duration string is unparseable, returns an error", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "cluster_stats_refresh_interval: not-a-duration\n")

		_, err := config.LoadFrom(path, "")

		r.Error(err)
	})

	t.Run("When the file exists but omits the key, falls back to the default interval", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "prometheus_url: http://prom.example.com:9090\n")

		cfg, err := config.LoadFrom(path, "")

		r.NoError(err)
		r.Equal(config.DefaultClusterStatsRefreshInterval, cfg.ClusterStatsRefreshInterval)
		r.Equal("http://prom.example.com:9090", cfg.PrometheusURL)
	})

	t.Run("When prometheus_url is absent, it stays empty", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "cluster_stats_refresh_interval: 5s\n")

		cfg, err := config.LoadFrom(path, "")

		r.NoError(err)
		r.Empty(cfg.PrometheusURL)
	})

	t.Run("When prometheus_placement is nomad, prometheus_url is not required", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "prometheus_placement: nomad\n")

		cfg, err := config.LoadFrom(path, "")

		r.NoError(err)
		r.Equal(config.PrometheusPlacementNomad, cfg.PrometheusPlacement)
	})

	t.Run("When prometheus_placement is nomad, custom job/port label override the defaults", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "prometheus_placement: nomad\nprometheus_job: my-prometheus\nprometheus_port_label: metrics\n")

		cfg, err := config.LoadFrom(path, "")

		r.NoError(err)
		r.Equal("my-prometheus", cfg.PrometheusJob)
		r.Equal("metrics", cfg.PrometheusPortLabel)
	})

	t.Run("When prometheus_placement is external and prometheus_url is set, it's valid", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "prometheus_placement: external\nprometheus_url: http://prom.example.com:9090\n")

		cfg, err := config.LoadFrom(path, "")

		r.NoError(err)
		r.Equal(config.PrometheusPlacementExternal, cfg.PrometheusPlacement)
		r.Equal("http://prom.example.com:9090", cfg.PrometheusURL)
	})

	t.Run("When prometheus_placement is external and prometheus_url is missing, returns an error", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "prometheus_placement: external\n")

		_, err := config.LoadFrom(path, "")

		r.Error(err)
	})

	t.Run("When prometheus_placement is invalid, returns an error", func(t *testing.T) {
		dir := t.TempDir()
		path := writeFile(t, dir, "prometheus_placement: bogus\n")

		_, err := config.LoadFrom(path, "")

		r.Error(err)
	})
}
