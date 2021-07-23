package nomad

import (
	"context"

	"github.com/hashicorp/nomad/api"
)

//go:generate counterfeiter . Client
type Client interface {
	Address() string
}

//go:generate counterfeiter . JobClient
type JobClient interface {
	List(*api.QueryOptions) ([]*api.JobListStub, *api.QueryMeta, error)
	Info(string, *api.QueryOptions) (*api.Job, *api.QueryMeta, error)
	Summary(string, *api.QueryOptions) (*api.JobSummary, *api.QueryMeta, error)
	Allocations(string, bool, *api.QueryOptions) ([]*api.AllocationListStub, *api.QueryMeta, error)
	Deregister(jobID string, purge bool, q *api.WriteOptions) (string, *api.WriteMeta, error)
	Register(job *api.Job, q *api.WriteOptions) (*api.JobRegisterResponse, *api.WriteMeta, error)
}

//go:generate counterfeiter . AllocationsClient
type AllocationsClient interface {
	List(*api.QueryOptions) ([]*api.AllocationListStub, *api.QueryMeta, error)
}

//go:generate counterfeiter . AllocFSClient
type AllocFSClient interface {
	Logs(alloc *api.Allocation, follow bool, task string, logType string, origin string, offset int64, cancel <-chan struct{}, q *api.QueryOptions) (<-chan *api.StreamFrame, <-chan error)
}

//go:generate counterfeiter . NamespaceClient
type NamespaceClient interface {
	List(*api.QueryOptions) ([]*api.Namespace, *api.QueryMeta, error)
}

//go:generate counterfeiter . DeploymentClient
type DeploymentClient interface {
	List(*api.QueryOptions) ([]*api.Deployment, *api.QueryMeta, error)
}

//go:generate counterfeiter . EventsClient
type EventsClient interface {
	Stream(ctx context.Context, topics map[api.Topic][]string, index uint64, q *api.QueryOptions) (<-chan *api.Events, error)
}

type SearchOptions struct {
	Namespace string
	Region    string
	DC        string
}

type Nomad struct {
	Client        Client
	EventsClient  EventsClient
	JobClient     JobClient
	NsClient      NamespaceClient
	AllocClient   AllocationsClient
	AllocFSClient AllocFSClient
	DpClient      DeploymentClient
}

func New(opts ...func(*Nomad) error) (*Nomad, error) {
	nomad := Nomad{}
	for _, opt := range opts {
		err := opt(&nomad)
		if err != nil {
			return nil, err
		}
	}

	return &nomad, nil
}

func Default(n *Nomad) error {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}
	// client.Allocations().Info(allocID string, q *api.QueryOptions)

	n.Client = client
	n.EventsClient = client.EventStream()
	n.JobClient = client.Jobs()
	n.NsClient = client.Namespaces()
	n.AllocClient = client.Allocations()
	n.AllocFSClient = client.AllocFS()
	n.DpClient = client.Deployments()

	return nil
}

func (n *Nomad) Address() string {
	return n.Client.Address()
}
