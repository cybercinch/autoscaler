// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import (
	"context"
	"net/http"
	"sync"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/linode/linodego"

	"golang.org/x/oauth2"
)

// provider implements a Linode provider.
type provider struct {
	init          sync.Once
	region        string
	token         string
	name          string
	instance_type string
	image         string
	ssh_key       string
	root_pass     string
	stackscript   string
	privateIP     bool
	userdata      *template.Template
	tags          []string
}

// New returns a new Linode provider.
func New(opts ...Option) autoscaler.Provider {
	p := new(provider)
	for _, opt := range opts {
		opt(p)
	}
	if p.region == "" {
		p.region = "ap-southeast"
	}
	if p.instance_type == "" {
		p.instance_type = "g6-standard-1"
	}
	if p.image == "" {
		p.image = "linode/ubuntu20.04"
	}
	if p.stackscript == "" {
		p.stackscript = "993696"
	}
	if p.userdata == nil {
		p.userdata = userdataT
	}
	return p
}

// helper function returns a new linode client.
func newClient(ctx context.Context, apiKey string) *linodego.Client {

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	linodeClient := linodego.NewClient(oauth2Client)
	linodeClient.SetDebug(true)

	return &linodeClient
}
