package models

import (
	"time"

	"github.com/hashicorp/nomad/api"
)

type HandlerFunc func(format string, args ...interface{})

type Handler string

const (
	HandleError Handler = Handler("Error")
	HandleFatal Handler = Handler("Fatal")
	HandleInfo  Handler = Handler("Info")

	TopicNamespace api.Topic = api.Topic("Namespace")
	TopicTaskGroup api.Topic = api.Topic("TaskGroup")
	TopicJobStatus api.Topic = api.Topic("JobStatus")
	TopicLog       api.Topic = api.Topic("Log")
)

type Job struct {
	ID                string
	Name              string
	Namespace         string
	Type              string
	Status            string
	StatusDescription string
	StatusSummary     Summary
	SubmitTime        time.Time
}

type JobStatus struct {
	ID                string
	Name              string
	Namespace         string
	Type              string
	Status            string
	StatusDescription string
	SubmitDate        time.Time
	Priority          int
	Datacenters       string
	Periodic          bool
	Parameterized     bool
	TaskGroups        []*TaskGroup
	TaskGroupStatus   []*TaskGroupStatus
	Allocations       []*Alloc
}

type Summary struct {
	Total   int
	Running int
}

type TaskGroup struct {
	Name     string
	JobID    string
	Queued   int
	Complete int
	Failed   int
	Running  int
	Starting int
	Lost     int
}

type TaskGroupStatus struct {
	ID                string
	Desired           int
	Placed            int
	Healthy           int
	Unhealthy         int
	ProgressDeadline  time.Duration
	Status            string
	StatusDescription string
}

type Alloc struct {
	ID            string
	Name          string
	Namespace     string
	TaskGroup     string
	Tasks         []AllocTask
	TaskNames     []string
	JobID         string
	JobType       string
	NodeID        string
	NodeName      string
	DesiredStatus string
	Version       uint64
	Status        string
	Created       time.Time
	Modified      time.Time
}

type AllocTask struct {
	Name   string
	Events []*api.TaskEvent
}

type Task struct {
	Name     string
	Image    string
	CPU      int
	MemoryMB int
	DiskMB   int
}

type Namespace struct {
	Name        string
	Description string
}

type Deployment struct {
	ID                string
	JobID             string
	Namespace         string
	Status            string
	StatusDescription string
}

type SearchResult struct {
}

type Status string

const (
	DesiredStatusRun  = "run"
	DesiredStatusStop = "stop"
	StatusRunning     = "running"
	StatusPending     = "pending"
	StatusDead        = "dead"
	StatusFailed      = "failed"
	StatusSuccessful  = "successful"

	TypeBatch   = "batch"
	TypeService = "service"
)

type Sentinel string

func (s Sentinel) Error() string {
	return string(s)
}

/*
Task Group  Desired  Placed  Healthy  Unhealthy  Progress Deadline
cadence     1        1       1        0          2021-05-12T14:27:03+02:00
iam         1        1       1        0          2021-05-12T14:28:37+02:00
*/
