// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import (
	"context"
	"errors"

	"github.com/drone/autoscaler/logger"
	"golang.org/x/sync/errgroup"
)

func (p *provider) setup(ctx context.Context) error {
	var g errgroup.Group
	if p.ssh_key == "" {
		g.Go(func() error {
			return p.setupKeypair(ctx)
		})
	}
	return g.Wait()
}

func (p *provider) setupKeypair(ctx context.Context) error {
	logger := logger.FromContext(ctx)

	logger.Debugln("finding default ssh key")

	client := newClient(ctx, p.token)
	res, err := client.ListSSHKeys(ctx, nil)
	if err != nil {
		return err
	}

	index := map[string]string{}
	for key := range res {
		logger.Debug(key)
		index[res[key].Label] = res[key].SSHKey
	}

	// if the account has multiple keys configured we will
	// attempt to use an existing key based on naming convention.
	for _, name := range []string{"drone", "id_rsa_drone"} {
		fingerprint, ok := index[name]
		if !ok {
			continue
		}
		p.ssh_key = fingerprint

		logger.
			WithField("name", name).
			WithField("fingerprint", fingerprint).
			Debugln("using matched ssh key")
		return nil
	}

	// if there were no matches but the account has at least
	// one keypair already created we will select the first
	// in the list.
	if len(res) > 0 {
		key := res[0].Label
		p.ssh_key = res[0].SSHKey

		logger.
			WithField("name", key).
			WithField("fingerprint", res[0].SSHKey).
			Debugln("using default ssh key")
		return nil
	}

	return errors.New("No matching keys")
}
