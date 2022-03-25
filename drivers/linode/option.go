// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

// Option configures a Linode provider option.
type Option func(*provider)

// WithImage returns an option to set the image.
func WithImage(image string) Option {
	return func(p *provider) {
		p.image = image
	}
}

// WithRegion returns an option to set the target region.
func WithRegion(region string) Option {
	return func(p *provider) {
		p.region = region
	}
}

// WithSize returns an option to set the instance size.
func WithType(instance_type string) Option {
	return func(p *provider) {
		p.instance_type = instance_type
	}
}

// WithSSHKey returns an option to set the ssh key.
func WithSSHKey(key string) Option {
	return func(p *provider) {
		p.ssh_key = key
	}
}

// WithTags returns an option to set the tags.
func WithTags(tags ...string) Option {
	return func(p *provider) {
		p.tags = tags
	}
}

// WithToken returns an option to set the auth token.
func WithToken(token string) Option {
	return func(p *provider) {
		p.token = token
	}
}

// WithPrivateIP returns an option to set the private IP address.
func WithPrivateIP(private bool) Option {
	return func(p *provider) {
		p.privateIP = private
	}
}

func WithRootPass(root_pass string) Option {
	return func(p *provider) {
		p.root_pass = root_pass
	}
}
