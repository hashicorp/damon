package watcher

import "github.com/hcjulz/damon/models"

// SubscribeToLogs starts an event stream for Logs
// which updates the state whenever a new log is written.
// The stream will be stopped whenever a new subscription happens.
func (w *Watcher) SubscribeToLogs(allocID, taskName, source string, notify func()) {
	// wipe any previous logs
	w.state.Logs = nil
	w.logResumer = &logResumer{
		allocID:  allocID,
		taskName: taskName,
		source:   source,
		notify:   notify,
	}

	alloc, ok := w.getAllocation(allocID)
	if !ok {
		w.NotifyHandler(models.HandleError, "allocation not found: %s", allocID)
		return
	}

	if len(alloc.TaskNames) == 0 {
		w.NotifyHandler(models.HandleError, "no tasks for allocation: %s", allocID)
		return
	}

	w.Subscribe(notify, models.TopicLog)
	w.Notify(models.TopicLog)

	cancel := make(chan struct{})
	streamCh, errorCh := w.nomad.Logs(allocID, taskName, source, cancel)

	w.activities.Add(cancel)

	go func() {
		for {
			select {
			case frame := <-streamCh:
				if frame.Data != nil {
					w.state.Logs = frame.Data
					w.Notify(models.TopicLog)
				}
			case err := <-errorCh:
				w.NotifyHandler(models.HandleError, err.Error())
			case <-cancel:
				return
			}
		}
	}()
}

func (w *Watcher) ResumeLogs() {
	w.SubscribeToLogs(w.logResumer.allocID, w.logResumer.taskName, w.logResumer.source, w.logResumer.notify)
}

func (w *Watcher) getAllocation(id string) (*models.Alloc, bool) {
	for _, a := range w.state.Allocations {
		if a.ID == id {
			return a, true
		}
	}

	return nil, false
}
