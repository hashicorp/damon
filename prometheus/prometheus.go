// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

// Package prometheus is an optional, standalone integration that
// queries an external Prometheus server for live cluster utilization
// metrics. It is deliberately independent of the nomad package: it
// talks to a different system (a Prometheus server that scrapes Nomad
// agents), not to Nomad's own HTTP API, and is only used when the
// user opts in via config.
package prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hcjulz/damon/models"
)

// Default PromQL queries for cluster-average host utilization. These
// rely on nomad.client.host.cpu.total_percent /
// nomad.client.host.memory.{used,total} being published, which
// requires `telemetry.publish_node_metrics = true` on the Nomad
// agents. There's no `..._used_percent` gauge for memory (confirmed
// against a live Nomad 2.0.4 server's metric names) — only raw
// used/total, so the percentage is computed here. Exact metric names
// can vary with the Prometheus sink's naming/relabeling config, so
// treat these as a best-effort default.
const (
	queryClusterCPUPercent    = "avg(nomad_client_host_cpu_total_percent)"
	queryClusterMemoryPercent = "100 * avg(nomad_client_host_memory_used) / avg(nomad_client_host_memory_total)"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type queryResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value [2]interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

// Query runs a PromQL instant query and returns its single scalar
// result. An empty result vector, a non-success status, or a
// non-2xx HTTP response are all treated as errors.
func (c *Client) Query(query string) (float64, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return 0, err
	}
	u.Path = "/api/v1/query"
	u.RawQuery = url.Values{"query": {query}}.Encode()

	resp, err := c.httpClient.Get(u.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("prometheus query failed: unexpected status %d", resp.StatusCode)
	}

	var qr queryResponse
	if err := json.NewDecoder(resp.Body).Decode(&qr); err != nil {
		return 0, err
	}

	if qr.Status != "success" {
		return 0, fmt.Errorf("prometheus query failed: status %q", qr.Status)
	}

	if len(qr.Data.Result) == 0 {
		return 0, fmt.Errorf("prometheus query returned no results for %q", query)
	}

	valueStr, ok := qr.Data.Result[0].Value[1].(string)
	if !ok {
		return 0, fmt.Errorf("prometheus query returned an unexpected value type for %q", query)
	}

	return strconv.ParseFloat(valueStr, 64)
}

// ClusterUtilization returns cluster-average CPU and memory
// utilization percentages.
func (c *Client) ClusterUtilization() (*models.ResourceUtilization, error) {
	cpu, err := c.Query(queryClusterCPUPercent)
	if err != nil {
		return nil, err
	}

	mem, err := c.Query(queryClusterMemoryPercent)
	if err != nil {
		return nil, err
	}

	return &models.ResourceUtilization{
		CPUPercent:    cpu,
		MemoryPercent: mem,
	}, nil
}
