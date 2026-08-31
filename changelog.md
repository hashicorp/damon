# Changelog

## 0.2.0

- Replaced the static ASCII-art logo in the header with a live cluster
  stats panel: node/server counts, leader address, job/alloc counts,
  and scheduled CPU/memory capacity vs. allocation. Refreshes on a
  configurable interval (default 10s) via a new `damon.yaml` config
  file (checked in the current directory, then
  `~/.config/damon/damon.yaml`).
- Added optional live cluster CPU/memory utilization, controlled by
  `prometheus_placement` in the config file: `none` (default),
  `external` (fixed `prometheus_url`), or `nomad` (auto-discover the
  address from a running Nomad job, re-resolved on every poll so it
  self-heals if the job moves — configurable via `prometheus_job` and
  `prometheus_port_label`, defaulting to `prometheus`/`http`).

## 0.1.0

- Initial tracked version.
