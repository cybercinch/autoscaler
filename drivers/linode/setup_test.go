// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import (
	"context"
	"testing"

	"github.com/h2non/gock"
)

func TestSetupKey_Single(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.linode.com").
		Get("/v4/profile/sshkeys").
		Reply(200).
		BodyString(respSingleKey)

	p := New(
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.ssh_key, "ssh-rsa 12345"; got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}
}

func TestSetupKey_FoundMatch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.linode.com").
		Get("/v4/profile/sshkeys").
		Reply(200).
		BodyString(respMultiKey)

	p := New(
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.ssh_key, "ssh-rsa 12346"; got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}
}

func TestSetupKey_NoMatch(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.linode.com").
		Get("/v4/profile/sshkeys").
		Reply(200).
		BodyString(respMultiKeyNoMatch)

	p := New(
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)

	err := p.setup(context.TODO())
	if err != nil {
		t.Error(err)
	}

	if got, want := p.ssh_key, "ssh-rsa 12347"; got != want {
		t.Errorf("Want fingerprint %s, got %s", want, got)
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}
}

var respSingleKey = `
{
	"data": [
	   {
		  "id": 147972,
		  "label": "aaron",
		  "ssh_key": "ssh-rsa 12345",
		  "created": "2021-10-05T21:57:05"
	   }
	],
	"page": 1,
	"pages": 1,
	"results": 1
}
`

var respMultiKey = `
{
	"data": [
	   {
			"id": 147972,
			"label": "aaron",
			"ssh_key": "ssh-rsa 12345",
			"created": "2021-10-05T21:57:05"
	   },
	   {
			"id": 147972,
			"label": "drone",
			"ssh_key": "ssh-rsa 12346",
			"created": "2021-10-05T21:57:05"
	   }
	],
	"page": 1,
	"pages": 1,
	"results": 2
}
`

var respMultiKeyNoMatch = `
{
	"data": [
	   {
			"id": 147972,
			"label": "My SSH Key 1",
			"ssh_key": "ssh-rsa 12347",
			"created": "2021-10-05T21:57:05"
	   },
	   {
			"id": 147972,
			"label": "My SSH Key 2",
			"ssh_key": "ssh-rsa 12348",
			"created": "2021-10-05T21:57:05"
	   }
	],
	"page": 1,
	"pages": 1,
	"results": 2
}
`
