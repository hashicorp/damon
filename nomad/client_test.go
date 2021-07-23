package nomad_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
)

func TestAddress(t *testing.T) {
	r := require.New(t)
	fakeClient := &nomadfakes.FakeClient{}
	nomad := &Nomad{Client: fakeClient}

	fakeClient.AddressReturns("127.0.0.1")

	addr := nomad.Address()

	r.Equal(addr, "127.0.0.1")
}
