package watcher

type ActivityPool struct {
	Activities []chan struct{}
}

func (a *ActivityPool) Add(act chan struct{}) {
	a.Activities = append(a.Activities, act)
}

func (a *ActivityPool) DeactivateAll() {
	for a.hasActivities() {
		a.deactivate()
	}
}

func (a *ActivityPool) deactivate() {
	if a.hasActivities() {
		ch := a.Activities[0]
		ch <- struct{}{}
		a.Activities = a.Activities[1:]
	}
}

func (a *ActivityPool) hasActivities() bool {
	return len(a.Activities) > 0
}
