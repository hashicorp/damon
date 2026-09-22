// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package refresher_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/refresher"
	"github.com/hcjulz/damon/refresher/refresherfakes"
)

func TestWatch(t *testing.T) {
	r := require.New(t)

	refreshFunc := &refresherfakes.FakeRefreshFunc{}

	fakeActivityPool := &refresherfakes.FakeActivities{}

	refresher := refresher.New(time.Second * 1).WithCustomActivityPool(fakeActivityPool)

	done := make(chan struct{})
	go func() {
		refresher.Refresh(refreshFunc.Spy)
		close(done)
	}()

	r.Eventually(func() bool {
		return refreshFunc.CallCount() > 1
	}, time.Second*6, time.Millisecond*100)

	stopChan := fakeActivityPool.AddArgsForCall(0)
	stopChan <- struct{}{}

	// the refresher stops on the signal, only then the channel is closed
	select {
	case <-done:
	case <-time.After(time.Second * 6):
		r.FailNow("refresh did not stop")
	}

	_, ok := <-stopChan
	r.False(ok)

	r.Equal(fakeActivityPool.DeactivateAllCallCount(), 1)
	r.Equal(fakeActivityPool.AddCallCount(), 1)
}
