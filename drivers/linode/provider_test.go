// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package linode

import "testing"

func TestDefaults(t *testing.T) {
	p := New().(*provider)
	if got, want := p.image, "private/15818922"; got != want {
		t.Errorf("Want image %q, got %q", want, got)
	}
	if got, want := p.region, "ap-southeast"; got != want {
		t.Errorf("Want region %q, got %q", want, got)
	}
	if got, want := p.instance_type, "g6-standard-1"; got != want {
		t.Errorf("Want size %q, got %q", want, got)
	}
	if got, want := p.ssh_key, ""; got != want {
		t.Errorf("Want key %q, got %q", want, got)
	}
	if got, want := p.token, ""; got != want {
		t.Errorf("Want token %q, got %q", want, got)
	}
	if got, want := len(p.tags), 0; got != want {
		t.Errorf("Want %d tags, got %d", want, got)
	}
}
