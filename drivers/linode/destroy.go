// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import (
	"context"
	"strconv"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
)

func (p *provider) Destroy(ctx context.Context, instance *autoscaler.Instance) error {
	logger := logger.FromContext(ctx).
		WithField("region", instance.Region).
		WithField("image", instance.Image).
		WithField("instance_type", instance.Size).
		WithField("label", instance.Name)

	client := newClient(ctx, p.token)
	id, err := strconv.Atoi(instance.ID)
	if err != nil {
		return err
	}

	_, err = client.GetInstance(ctx, id)
	if err != nil {
		logger.WithError(err).
			Warnln("linode does not exist")
		return err
	} else if err != nil {
		logger.WithError(err).
			Errorln("cannot find linode")
		return err
	}

	logger.Debugln("deleting linode")

	err = client.DeleteInstance(ctx, id)
	if err != nil {
		logger.WithError(err).
			Errorln("deleting linode failed")
		return err
	}

	logger.Debugln("linode deleted")

	return nil
}
