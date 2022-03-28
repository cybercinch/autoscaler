// Copyright 2020 Aaron Guise
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/logger"
	"github.com/linode/linodego"
)

func defValue(value *string, def string) string {
	var return_val = ""
	if value != nil {
		return_val = *value
	} else {
		return_val = def
	}
	return return_val
}

func (p *provider) Create(ctx context.Context, opts autoscaler.InstanceCreateOpts) (*autoscaler.Instance, error) {
	p.init.Do(func() {
		p.setup(ctx)
	})

	buf := new(bytes.Buffer)
	err := p.userdata.Execute(buf, &opts)
	if err != nil {
		return nil, err
	}

	name := strings.ToLower(opts.Name)

	logger := logger.FromContext(ctx).
		WithField("image", p.image).
		WithField("size", p.instance_type).
		WithField("label", name)

	logger.Debugln("instance insert")

	instance := &autoscaler.Instance{
		Provider: autoscaler.ProviderLinode,
		ID:       "",
		Name:     name,
		Image:    p.image,
		Region:   p.region,
		Size:     p.instance_type,
	}

	client := newClient(ctx, p.token)

	stack_id, err := strconv.ParseInt(p.stackscript, 0, 32)
	userdata_map := make(map[string]string)
	userdata_map["userdata"] = b64.StdEncoding.EncodeToString([]byte(buf.String()))

	linode, err := client.CreateInstance(ctx, linodego.InstanceCreateOptions{
		Region:          p.region,
		Type:            p.instance_type,
		Label:           name,
		Image:           p.image,
		StackScriptID:   int(stack_id),
		StackScriptData: userdata_map,
		AuthorizedKeys:  []string{p.ssh_key},
		RootPass:        p.root_pass,
	})

	if err != nil {
		logger.Errorln(err.Error())
		return instance, err
	}

	logger.Debug("Created Instance %v\n", linode)
	logger.Debug("Error is %v\n", err)
	if linode != nil {
		check_instance, err := client.GetInstance(ctx, linode.ID)
		logger.Debug("Error: %v\n\n", err)
		logger.Debug("Variable: %v\n\n", check_instance)
		if err != nil {
			logger.Errorln(err.Error())
			return instance, err
		}
		if check_instance != nil {
			logger.Debug("Polled Instance ID = %d\n", check_instance.ID)
			if check_instance.Status != "running" {
			poller:
				for {
					time.Sleep(3 * time.Second)
					logger.Debug("Instance ID = %d\n", linode.ID)
					poll_instance, err := client.GetInstance(ctx, linode.ID)
					if err != nil {
						logger.Errorln(err)
						return instance, err
					}
					logger.Debug("Polled Instance ID = %d\n", poll_instance.ID)
					logger.Debug("%v", poll_instance)
					if poll_instance.Status == "running" {
						// The server has finished provisioning
						break poller
					}

				}
			}
		}

	}

	instance_id := ""

	if linode != nil {
		instance_id = strconv.Itoa(linode.ID)
	}
	instance.ID = instance_id
	instance.Address = linode.IPv4[0].String()

	return instance, err
}
