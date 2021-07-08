package refresher

import (
	"time"

	"github.com/hcjulz/damon/watcher"
)

const defaultRefreshInterval = time.Second * 2

//go:generate counterfeiter . RefreshFunc
type RefreshFunc func()

//go:generate counterfeiter . Activities
type Activities interface {
	Add(chan struct{})
	DeactivateAll()
}

type Refresher struct {
	activities      Activities
	RefreshInterval time.Duration
}

func New(d time.Duration) *Refresher {
	if d == 0 {
		d = defaultRefreshInterval
	}

	return &Refresher{
		activities:      &watcher.ActivityPool{},
		RefreshInterval: d,
	}
}

func (w *Refresher) WithCustomActivityPool(a Activities) *Refresher {
	w.activities = a
	return w
}

func (w *Refresher) Refresh(refresh RefreshFunc) {
	stop := make(chan struct{}, 1)
	w.activities.DeactivateAll()
	w.activities.Add(stop)

	refresh()

	ticker := time.NewTicker(w.RefreshInterval)
	for range ticker.C {
		select {
		case <-stop:
			close(stop)
			return
		default:
			refresh()
		}
	}
}
