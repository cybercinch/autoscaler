// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import (
	"context"
	"testing"

	"github.com/drone/autoscaler"
	"github.com/linode/linodego"

	"github.com/h2non/gock"
)

func TestCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.linode.com").
		Post("v4/linode/instances").
		Reply(200).
		BodyString(respCreateInstance)

	gock.New("https://api.linode.com").
		Get("v4/linode/instances/35479107").
		Reply(200).
		BodyString(respDescInstance)

	p := New(
		WithSSHKey("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)
	p.init.Do(func() {}) // prevent init function

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err != nil {
		t.Error(err)
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}

	t.Run("Attributes", testInstance(instance))
	t.Run("Address", testInstanceAddress(instance))
}

func TestCreate_CreateError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.linode.com").
		Post("/v4/linode/instances").
		Reply(400).
		BodyString(respCreateInstanceError).
		Header.Set("Content-Type", "application/json")

	p := New(
		WithSSHKey("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)
	p.init.Do(func() {}) // prevent init function

	_, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expect error returned from digital ocean")
	} else if _, ok := err.(*linodego.Error); !ok {
		t.Errorf("Expect ErrorResponse digital ocean")
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}
}

func TestCreate_DescribeError(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.linode.com").
		Post("/v4/linode/instances").
		Reply(200).
		BodyString(respCreateInstance)

	gock.New("https://api.linode.com").
		Get("/v4/linode/instances/35479107").
		Reply(404).
		BodyString(respDescInstanceError).
		SetHeader("Content-Type", "application/json")

	p := New(
		WithSSHKey("58:8e:30:66:fc:e2:ff:ad:4f:6f:02:4b:af:28:0d:c7"),
		WithToken("77e027c7447f468068a7d4fea41e7149a75a94088082c66fcf555de3977f69d3"),
	).(*provider)
	p.init.Do(func() {}) // prevent init function

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent1"})
	if err == nil {
		t.Errorf("Expect error returned from Linode")
	} else if _, ok := err.(*linodego.Error); !ok {
		t.Errorf("Expect ErrorResponse digital ocean")
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}

	t.Run("Attributes", testInstanceFailed(instance))
}

func testInstance(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.ID, "35479107"; got != want {
			t.Errorf("Want Linode ID %v, got %v", want, got)
		}
		if got, want := instance.Image, "private/15818922"; got != want {
			t.Errorf("Want Linode Image %v, got %v", want, got)
		}
		if got, want := instance.Name, "agent1"; got != want {
			t.Errorf("Want Linode Name %v, got %v", want, got)
		}
		if got, want := instance.Region, "ap-southeast"; got != want {
			t.Errorf("Want Linode Region %v, got %v", want, got)
		}
		if got, want := instance.Provider, autoscaler.ProviderLinode; got != want {
			t.Errorf("Want Linode Provider %v, got %v", want, got)
		}
	}
}

func testInstanceFailed(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.ID, ""; got != want {
			t.Errorf("Want Linode ID %v, got %v", want, got)
		}
		if got, want := instance.Image, "private/15818922"; got != want {
			t.Errorf("Want Linode Image %v, got %v", want, got)
		}
		if got, want := instance.Name, "agent1"; got != want {
			t.Errorf("Want Linode Name %v, got %v", want, got)
		}
		if got, want := instance.Region, "ap-southeast"; got != want {
			t.Errorf("Want Linode Region %v, got %v", want, got)
		}
		if got, want := instance.Provider, autoscaler.ProviderLinode; got != want {
			t.Errorf("Want Linode Provider %v, got %v", want, got)
		}
	}
}

func testInstanceAddress(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.Address, "172.105.255.97"; got != want {
			t.Errorf("Want droplet Address %v, got %v", want, got)
		}
	}
}

// sample response for POST /v2/droplets
const respCreateInstance = `
{
	"id": 35479107,
	"label": "linode35479107",
	"group": "",
	"status": "provisioning",
	"created": "2022-03-21T08:41:55",
	"updated": "2022-03-21T08:41:55",
	"type": "g6-standard-2",
	"ipv4": [
	   "172.105.255.97"
	],
	"ipv6": "2400:8907::f03c:93ff:fe28:82f9/128",
	"image": null,
	"region": "ap-southeast",
	"specs": {
	   "disk": 81920,
	   "memory": 4096,
	   "vcpus": 2,
	   "gpus": 0,
	   "transfer": 4000
	},
	"alerts": {
	   "cpu": 180,
	   "network_in": 10,
	   "network_out": 10,
	   "transfer_quota": 80,
	   "io": 10000
	},
	"backups": {
	   "enabled": false,
	   "schedule": {
		  "day": null,
		  "window": null
	   },
	   "last_successful": null
	},
	"hypervisor": "kvm",
	"watchdog_enabled": true,
	"tags": []
 }
`

// sample response for POST /v2/droplets/:id
const respDescInstance = `
{
	"id": 35479107,
	"label": "linode35479107",
	"group": "",
	"status": "running",
	"created": "2022-03-21T08:41:55",
	"updated": "2022-03-21T08:41:55",
	"type": "g6-standard-2",
	"ipv4": [
	   "172.105.255.97"
	],
	"ipv6": "2400:8907::f03c:93ff:fe28:82f9/128",
	"image": null,
	"region": "ap-southeast",
	"specs": {
	   "disk": 81920,
	   "memory": 4096,
	   "vcpus": 2,
	   "gpus": 0,
	   "transfer": 4000
	},
	"alerts": {
	   "cpu": 180,
	   "network_in": 10,
	   "network_out": 10,
	   "transfer_quota": 80,
	   "io": 10000
	},
	"backups": {
	   "enabled": false,
	   "schedule": {
		  "day": null,
		  "window": null
	   },
	   "last_successful": null
	},
	"hypervisor": "kvm",
	"watchdog_enabled": true,
	"tags": []
 }
`
const respDescInstanceError = `
{
   "errors": [
      {
         "reason": "Not found"
      }
   ]
}
`
const respCreateInstanceError = `
{
	"errors": [
	   {
		  "reason": "A valid plan type by that ID was not found",
		  "field": "type"
	   }
	]
 }
`
