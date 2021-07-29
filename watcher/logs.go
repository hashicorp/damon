package watcher

import "github.com/hcjulz/damon/models"

// SubscribeToLogs starts an event stream for Logs
// which updates the state whenever a new log is written.
// The stream will be stopped whenever a new subscription happens.
func (w *Watcher) SubscribeToLogs(allocID, source string, notify func()) {
	// wipe any previous logs
	w.state.Logs = nil

	alloc, ok := w.getAllocation(allocID)
	if !ok {
		w.NotifyHandler(models.HandleError, "allocation not found: %s", allocID)
		return
	}

	if len(alloc.TaskNames) == 0 {
		w.NotifyHandler(models.HandleError, "no tasks for allocation: %s", allocID)
		return
	}

	w.Subscribe(models.TopicLog, notify)
	w.Notify(models.TopicLog)

	cancel := make(chan struct{})
	streamCh, errorCh := w.nomad.Logs(allocID, alloc.TaskNames[0], source, cancel)

	w.activities.Add(cancel)

	go func() {
		for {
			select {
			case frame := <-streamCh:
				w.state.Logs = append(w.state.Logs, frame.Data...)
				w.Notify(models.TopicLog)
			case err := <-errorCh:
				w.NotifyHandler(models.HandleError, err.Error())
			case <-cancel:
				return
			}
		}
	}()
}

func (w *Watcher) getAllocation(id string) (*models.Alloc, bool) {
	for _, a := range w.state.Allocations {
		if a.ID == id {
			return a, true
		}
	}

	return nil, false
}
