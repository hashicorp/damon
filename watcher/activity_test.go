package watcher_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/watcher"
)

func TestAdd(t *testing.T) {
	r := require.New(t)

	activity := &watcher.ActivityPool{}

	activity.Add(make(chan bool))
	r.Equal(len(activity.Activities), 1)

	activity.Add(make(chan bool))
	r.Equal(len(activity.Activities), 2)
}

func TestDeactivateAll(t *testing.T) {
	r := require.New(t)

	activity := &watcher.ActivityPool{}
	activity.Activities = []chan bool{
		make(chan bool, 1),
		make(chan bool, 1),
	}

	activity.DeactivateAll()

	r.Empty(activity.Activities)
}
