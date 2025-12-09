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

	go refresher.Refresh(refreshFunc.Spy)

	r.Eventually(func() bool {
		return refreshFunc.CallCount() > 1
	}, time.Second*6, time.Second*1)

	stopChan := fakeActivityPool.AddArgsForCall(0)
	stopChan <- struct{}{}

	// check if the channel is eventually closed
	r.Eventually(func() bool {
		_, ok := <-stopChan
		return !ok
	}, time.Second*6, time.Second*2)

	r.Equal(fakeActivityPool.DeactivateAllCallCount(), 1)
	r.Equal(fakeActivityPool.AddCallCount(), 1)
}
