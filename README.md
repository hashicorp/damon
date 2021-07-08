# Damon - A terminal Dashboard for HashiCorp Nomad

Damon is a terminal user interface (TUI) for Nomad. It provides functionality to observe and interact with Nomad resources such as Jobs, Deployments, or Allocations.

**Additional Notes**

Damon is in an early stage and is under active development. We are working on improving the performance and adding new features to Damon.
Please take a look at the Damon [project board](https://github.com/hashicorp/damon/projects/2) to see what features you can expect in near future.
If you find a bug or you have great ideas how Damon can be improved feel free to open an issue. To avoid duplicates, please check the [project board](https://github.com/hashicorp/damon/projects/2) 
before submitting one. Thank you!

## Installation

### Brew

--> Comming soon

### Building from source

Make sure you have your go environment setup and run:

1. Clone the project
1. Run `$ make build` 
1. Run `$ make run`

### How to use it

Once `Damon` is installed, simply run:

```
$ damon
```

#### Environment Variables

Damon reads the following environment variables on startup:

- `NOMAD_TOKEN`
- `NOMAD_ADDR`
- `NOMAD_REGION`
- `NOMAD_NAMESPACE`
- `NOMAD_HTTP_AUTH`
- `NOMAD_CACERT`
- `NOMAD_CAPATH`
- `NOMAD_CLIENT_CERT`
- `NOMAD_CLIENT_KEY`
- `NOMAD_TLS_SERVER_NAME`
- `NOMAD_SKIP_VERIFY`

You can read about them in detail [here](https://www.nomadproject.io/docs/runtime/environment).

## Navigation

### General

On every table or text view, you can use:

- `k` or `arrow up` to navigate up
- `j` or `arrow down` to navigate down

### Top Level Commands

- Show Jobs: `ctrl-j`
- Show Deployments: `ctrl-d`
- Show Namespaces: `ctrl-n`
- Jump to a Jobs Allocations: `ctrl-j`
- Switch Namespace: `s`
- Quit: `ctrl-c`

### Job View Commands

- Show Allocations for a Job: `<ENTER>` (on the selected job)
- Show TaskGroups for a Job: `<t>` (on the selected job)
- Filter Job: `</>` (on the selected job)

### Allocation View Commands

- Show logs on `STDOUT` for an allocation: `<ENTER>`
- Show logs on `STDERR` for an allocation: `<ctrl-e>`

### Log View

When Damon displays logs, you can navigate through the logs using `j`, `k`, `G`, and `g`.
To filter logs you can hit `/` which will open a input field to enter the filter string.
