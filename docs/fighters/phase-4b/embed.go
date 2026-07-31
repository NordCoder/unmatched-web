// Package phase4bfighters exposes the verified Phase 4B launch-fighter
// manifests to the playable-slice content registry.
package phase4bfighters

import _ "embed"

var (
	//go:embed robin-hood.yaml
	RobinHood []byte

	//go:embed bigfoot.yaml
	Bigfoot []byte
)
