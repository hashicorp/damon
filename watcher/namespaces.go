package watcher

import (
	"time"

	"github.com/hcjulz/damon/models"
)

// SubscribeToNamespaces starts a goroutine to poll Namespaces based
// on the provided interval. It updates the state accordingly.
// The goroutine will be stopped whenever a new subscription happens.
func (w *Watcher) SubscribeToNamespaces(notify func()) {
	w.updateNamespaces()
	w.Subscribe(models.TopicNamespace, notify)
	w.Notify(models.TopicNamespace)

	stop := make(chan struct{})
	w.activities.Add(stop)

	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.updateNamespaces()
				w.Notify(models.TopicNamespace)
			case <-stop:
				return
			}
		}
	}()
}

func (w *Watcher) updateNamespaces() {
	ns, err := w.nomad.Namespaces(nil)
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
	}

	w.state.Namespaces = ns
}
