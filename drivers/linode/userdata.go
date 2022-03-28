// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import "github.com/drone/autoscaler/drivers/internal/userdata"

var userdataT = userdata.Parse(`#cloud-config
mkdir -p /etc/systemd/system/docker.service.d
cat > /etc/systemd/system/docker.service.d/override.conf <<'EOS'
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd
EOS

cat > /etc/default/docker <<'EOS'
# Unset all docker options as configured in
# /etc/docker/daemon.json
DOCKER_OPTS=""
EOS

cat > /etc/docker/daemon.json <<'EOS'
{
  "dns": [ "8.8.8.8", "8.8.4.4" ],
  "hosts": [ "0.0.0.0:2376", "unix:///var/run/docker.sock" ],
  "tls": true,
  "tlsverify": true,
  "tlscacert": "/etc/docker/ca.pem",
  "tlscert": "/etc/docker/server-cert.pem",
  "tlskey": "/etc/docker/server-key.pem"
}
EOS

echo "{{ .CACert | base64 }}" | base64 --decode > /etc/docker/ca.pem
echo "{{ .TLSCert | base64 }}" | base64 --decode > /etc/docker/server-cert.pem
echo "{{ .TLSKey | base64 }}" | base64 --decode > /etc/docker/server-key.pem

systemctl daemon-reload
systemctl restart docker
`)
