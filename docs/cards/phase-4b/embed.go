// Package phase4bcards exposes the verified Phase 4B launch-deck manifests to
// the playable-slice content registry. Gameplay code must not interpret YAML.
package phase4bcards

import _ "embed"

var (
	//go:embed robin-hood.yaml
	RobinHood []byte

	//go:embed bigfoot.yaml
	Bigfoot []byte
)
