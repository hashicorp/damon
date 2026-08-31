// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package prometheus_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/prometheus"
)

func TestQuery(t *testing.T) {
	r := require.New(t)

	t.Run("When the query succeeds", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[{"metric":{},"value":[1234,"42.5"]}]}}`))
		}))
		defer server.Close()

		client := prometheus.New(server.URL)
		value, err := client.Query("avg(some_metric)")

		r.NoError(err)
		r.Equal(42.5, value)
	})

	t.Run("When the HTTP response is non-2xx", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := prometheus.New(server.URL)
		_, err := client.Query("avg(some_metric)")

		r.Error(err)
	})

	t.Run("When the response body is malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`not json`))
		}))
		defer server.Close()

		client := prometheus.New(server.URL)
		_, err := client.Query("avg(some_metric)")

		r.Error(err)
	})

	t.Run("When the status is not success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`{"status":"error","data":{}}`))
		}))
		defer server.Close()

		client := prometheus.New(server.URL)
		_, err := client.Query("avg(some_metric)")

		r.Error(err)
	})

	t.Run("When the result vector is empty", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Write([]byte(`{"status":"success","data":{"resultType":"vector","result":[]}}`))
		}))
		defer server.Close()

		client := prometheus.New(server.URL)
		_, err := client.Query("avg(some_metric)")

		r.Error(err)
	})
}

func TestClusterUtilization(t *testing.T) {
	r := require.New(t)

	t.Run("When both queries succeed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			query := req.URL.Query().Get("query")
			switch query {
			case "avg(nomad_client_host_cpu_total_percent)":
				w.Write([]byte(`{"status":"success","data":{"result":[{"value":[1234,"55.5"]}]}}`))
			case "100 * avg(nomad_client_host_memory_used) / avg(nomad_client_host_memory_total)":
				w.Write([]byte(`{"status":"success","data":{"result":[{"value":[1234,"70.1"]}]}}`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		client := prometheus.New(server.URL)
		util, err := client.ClusterUtilization()

		r.NoError(err)
		r.Equal(55.5, util.CPUPercent)
		r.Equal(70.1, util.MemoryPercent)
	})

	t.Run("When the CPU query fails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := prometheus.New(server.URL)
		util, err := client.ClusterUtilization()

		r.Error(err)
		r.Nil(util)
	})
}
