// Package transport contains transport-neutral delivery envelopes. Wave 1 does
// not select HTTP, WebSocket, or a public schema technology.
package transport

import (
	"encoding/json"

	"github.com/NordCoder/unmatched-web/internal/application"
)

type ProjectionEnvelope struct {
	SchemaVersion string                       `json:"schema_version"`
	Type          string                       `json:"type"`
	Projection    application.PlayerProjection `json:"projection"`
}

func EncodeProjection(projection application.PlayerProjection) ([]byte, error) {
	return json.Marshal(ProjectionEnvelope{
		SchemaVersion: "core-projection/v1",
		Type:          "player_projection",
		Projection:    projection,
	})
}
