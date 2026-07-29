// Package phase4bcards exposes the verified Phase 4B launch-deck manifests to
// the playable-slice content registry. Gameplay code must not interpret YAML.
package phase4bcards

import (
	"bytes"
	_ "embed"
)

var (
	//go:embed robin-hood.yaml
	robinHoodManifest []byte

	//go:embed bigfoot.yaml
	bigfootManifest []byte

	RobinHood = runtimeManifest(robinHoodManifest)
	Bigfoot   = runtimeManifest(bigfootManifest)
)

// runtimeManifest retains manifest metadata and the complete cards section but
// excludes later top-level evidence and validation lists. Those records remain
// in the source YAML and are validated independently; they are not game cards.
func runtimeManifest(raw []byte) []byte {
	lines := bytes.SplitAfter(raw, []byte("\n"))
	seenCards := false
	end := 0
	for _, line := range lines {
		trimmed := bytes.TrimRight(line, "\r\n")
		if bytes.Equal(trimmed, []byte("cards:")) {
			seenCards = true
		} else if seenCards && len(trimmed) > 0 && trimmed[0] != ' ' {
			break
		}
		end += len(line)
	}
	if !seenCards {
		return append([]byte(nil), raw...)
	}
	return append([]byte(nil), raw[:end]...)
}
