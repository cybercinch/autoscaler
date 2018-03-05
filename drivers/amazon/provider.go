// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package amazon

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/drone/autoscaler"
)

const (
	maxRetries = 50

	defaultDeviceName = "/dev/sda1"
	defaultImage      = "ami-66506c1c"
	defaultRootSize   = 16
	defaultVolumeType = "gp2"
)

type provider struct {
	key    string
	region string
	image  string
	size   string
	subnet string
	groups []string
	tags   []string
}

func (p *provider) getClient() *ec2.EC2 {
	config := aws.NewConfig()
	config = config.WithRegion(p.region)
	config = config.WithMaxRetries(maxRetries)
	return ec2.New(session.New(config))
}

// New returns a new Digital Ocean provider.
func New(opts ...Option) autoscaler.Provider {
	p := &provider{
		region: "us-east-1",
		size:   "t2.medium",
		image:  "ami-66506c1c",
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Env returns true if the Digital Ocean provider
// environment variables are set.
func Env() bool {
	return os.Getenv("AWS_ACCESS_KEY_ID") != "" &&
		os.Getenv("AWS_SECRET_ACCESS_KEY") != ""
}
