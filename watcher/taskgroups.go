// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package watcher

import (
	"time"

	"github.com/hcjulz/damon/models"
)

// SubscribeToTaskGroups starts a goroutine to polls TaskGroups every two
// seconds to update the state. The goroutine will be stopped whenever
// a new subscription happens.
func (w *Watcher) SubscribeToTaskGroups(jobID string, notify func()) error {
	w.updateTaskGroups(jobID)
	w.Subscribe(notify, models.TopicTaskGroup)
	w.Notify(models.TopicTaskGroup)

	stop := make(chan struct{})
	w.activities.Add(stop)

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.updateTaskGroups(jobID)
				w.Notify(models.TopicTaskGroup)
			case <-stop:
				return
			}
		}
	}()

	return nil
}

func (w *Watcher) updateTaskGroups(jobID string) {
	tg, err := w.nomad.TaskGroups(jobID, nil)
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
	}

	w.state.TaskGroups = tg
}
