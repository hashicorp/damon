// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package nomad_test

import (
	"errors"
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/require"

	. "github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/nomad/nomadfakes"
)

func TestNamespaces(t *testing.T) {
	r := require.New(t)
	fakeNsClient := &nomadfakes.FakeNamespaceClient{}
	nomad := &Nomad{NsClient: fakeNsClient}

	t.Run("When everything is fine", func(t *testing.T) {
		fakeNsClient.ListReturns([]*api.Namespace{
			{
				Name:        "default",
				Description: "the default namespace",
			},
			{
				Name:        "test",
				Description: "the test namespace",
			},
		}, nil, nil)

		ns, err := nomad.Namespaces(nil)

		r.NoError(err)
		r.Len(ns, 2)

		r.Equal(ns[0].Name, "default")
		r.Equal(ns[0].Description, "the default namespace")

		r.Equal(ns[1].Name, "test")
		r.Equal(ns[1].Description, "the test namespace")
	})

	t.Run("When everything is fine", func(t *testing.T) {
		fakeNsClient.ListReturns(nil, nil, errors.New("fail!"))

		ns, err := nomad.Namespaces(nil)

		r.Nil(ns)
		r.Error(err)
		r.EqualError(err, "fail!")
	})
}
