package watcher

import (
	"time"

	"github.com/hcjulz/damon/models"
)

// SubscribeToJobStatus starts a goroutine to polls JobStatus every two
// seconds to update the state. The goroutine will be stopped whenever
// a new subscription happens.
func (w *Watcher) SubscribeToJobStatus(jobID string, notify func()) error {
	w.updateJobStatus(jobID)
	w.Subscribe(models.TopicJobStatus, notify)
	w.Notify(models.TopicJobStatus)

	stop := make(chan struct{})
	w.activities.Add(stop)

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.updateJobStatus(jobID)
				w.Notify(models.TopicJobStatus)
			case <-stop:
				return
			}
		}
	}()

	return nil
}

func (w *Watcher) updateJobStatus(jobID string) {
	js, err := w.nomad.JobStatus(jobID, nil)
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
	}

	w.state.JobStatus = js
}
