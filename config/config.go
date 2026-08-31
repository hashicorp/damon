// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	FileName                           = "damon.yaml"
	DefaultClusterStatsRefreshInterval = 10 * time.Second

	DefaultPrometheusJob       = "prometheus"
	DefaultPrometheusPortLabel = "http"
)

// PrometheusPlacement selects how (or whether) damon locates a
// Prometheus server for live cluster utilization metrics.
type PrometheusPlacement string

const (
	// PrometheusPlacementNone disables the Prometheus integration.
	PrometheusPlacementNone PrometheusPlacement = "none"
	// PrometheusPlacementNomad discovers Prometheus's address from a
	// running Nomad job/allocation on each poll.
	PrometheusPlacementNomad PrometheusPlacement = "nomad"
	// PrometheusPlacementExternal uses a fixed, externally-configured
	// Prometheus URL.
	PrometheusPlacementExternal PrometheusPlacement = "external"
)

type Config struct {
	ClusterStatsRefreshInterval time.Duration
	PrometheusPlacement         PrometheusPlacement
	PrometheusURL               string
	PrometheusJob               string
	PrometheusPortLabel         string
}

type fileConfig struct {
	ClusterStatsRefreshInterval string `yaml:"cluster_stats_refresh_interval"`
	PrometheusPlacement         string `yaml:"prometheus_placement"`
	PrometheusURL               string `yaml:"prometheus_url"`
	PrometheusJob               string `yaml:"prometheus_job"`
	PrometheusPortLabel         string `yaml:"prometheus_port_label"`
}

// Load looks for ./damon.yaml (the current working directory) first,
// then falls back to ~/.config/damon/damon.yaml. If neither exists,
// it returns defaults. If a file exists but is malformed, it returns
// an error so callers can fail fast rather than silently ignore a typo.
func Load() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	var homePath string
	if home != "" {
		homePath = filepath.Join(home, ".config", "damon", FileName)
	}

	return LoadFrom(filepath.Join(cwd, FileName), homePath)
}

// LoadFrom loads config from the first of cwdPath/homePath that
// exists, in that order. Either path may be empty to skip it.
func LoadFrom(cwdPath, homePath string) (*Config, error) {
	for _, path := range []string{cwdPath, homePath} {
		if path == "" {
			continue
		}

		fc, found, err := load(path)
		if err != nil {
			return nil, err
		}
		if found {
			return toConfig(fc)
		}
	}

	return toConfig(&fileConfig{})
}

func load(path string) (*fileConfig, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var fc fileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return nil, false, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return &fc, true, nil
}

func toConfig(fc *fileConfig) (*Config, error) {
	cfg := &Config{
		ClusterStatsRefreshInterval: DefaultClusterStatsRefreshInterval,
		PrometheusPlacement:         PrometheusPlacementNone,
		PrometheusURL:               fc.PrometheusURL,
		PrometheusJob:               DefaultPrometheusJob,
		PrometheusPortLabel:         DefaultPrometheusPortLabel,
	}

	if fc.ClusterStatsRefreshInterval != "" {
		d, err := time.ParseDuration(fc.ClusterStatsRefreshInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid cluster_stats_refresh_interval %q: %w", fc.ClusterStatsRefreshInterval, err)
		}
		cfg.ClusterStatsRefreshInterval = d
	}

	if fc.PrometheusJob != "" {
		cfg.PrometheusJob = fc.PrometheusJob
	}

	if fc.PrometheusPortLabel != "" {
		cfg.PrometheusPortLabel = fc.PrometheusPortLabel
	}

	if fc.PrometheusPlacement != "" {
		cfg.PrometheusPlacement = PrometheusPlacement(fc.PrometheusPlacement)
	}

	switch cfg.PrometheusPlacement {
	case PrometheusPlacementNone, PrometheusPlacementNomad:
		// no further validation needed
	case PrometheusPlacementExternal:
		if cfg.PrometheusURL == "" {
			return nil, fmt.Errorf("prometheus_placement: %q requires prometheus_url to be set", PrometheusPlacementExternal)
		}
	default:
		return nil, fmt.Errorf(
			"invalid prometheus_placement %q: must be %q, %q, or %q",
			cfg.PrometheusPlacement, PrometheusPlacementNone, PrometheusPlacementNomad, PrometheusPlacementExternal,
		)
	}

	return cfg, nil
}
