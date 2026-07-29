package application

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
)

func TestCanonicalCommandIdentityNormalizesDefaultsOrderAndNFC(t *testing.T) {
	revision := uint64(7)
	first := contracts.Command{
		ID:               "command-1",
		SchemaVersion:    CommandSchemaV1,
		Type:             CommandStartAction,
		MatchID:          "match-1",
		ActorPlayerID:    "player-1",
		ExpectedRevision: &revision,
		Payload:          json.RawMessage(`{"kind":"MANEUVER","context":{"label":"e\u0301","b":2,"a":1}}`),
	}
	second := first
	second.Payload = json.RawMessage(`{"source_ref":"","context":{"a":1,"b":2,"label":"é"},"kind":"MANEUVER"}`)

	_, normalizedFirst, firstFingerprint, err := normalizeCommandRequest("principal-1", first)
	if err != nil {
		t.Fatalf("normalize first command: %v", err)
	}
	_, normalizedSecond, secondFingerprint, err := normalizeCommandRequest("principal-1", second)
	if err != nil {
		t.Fatalf("normalize second command: %v", err)
	}
	if !bytes.Equal(firstFingerprint, secondFingerprint) {
		t.Fatalf("equivalent commands produced different fingerprints: %x != %x", firstFingerprint, secondFingerprint)
	}
	if !bytes.Equal(normalizedFirst.Payload, normalizedSecond.Payload) {
		t.Fatalf("equivalent commands produced different payloads: %s != %s", normalizedFirst.Payload, normalizedSecond.Payload)
	}
	if got, want := string(normalizedFirst.Payload), `{"context":{"a":1,"b":2,"label":"é"},"kind":"MANEUVER","source_ref":""}`; got != want {
		t.Fatalf("unexpected normalized payload: got %s want %s", got, want)
	}
}

func TestCanonicalJSONUsesRFC8785StringAndUTF16Ordering(t *testing.T) {
	value, err := parseCanonicalJSON([]byte(`{"\ue000":2,"\ud800\udc00":1,"text":"<>&\n"}`))
	if err != nil {
		t.Fatalf("parse canonical vector: %v", err)
	}
	if got, want := string(canonicalJSON(value)), "{\"text\":\"<>&\\n\",\"𐀀\":1,\"\":2}"; got != want {
		t.Fatalf("unexpected canonical bytes: got %q want %q", got, want)
	}
}

func TestCanonicalCommandIdentityRejectsClosedSchemaViolations(t *testing.T) {
	revision := uint64(1)
	tests := []struct {
		name    string
		payload string
	}{
		{name: "unknown field", payload: `{"kind":"MANEUVER","unknown":true}`},
		{name: "duplicate field", payload: `{"kind":"MANEUVER","kind":"ATTACK"}`},
		{name: "NFC duplicate key", payload: `{"kind":"MANEUVER","context":{"e\u0301":1,"é":2}}`},
		{name: "unsafe integer", payload: `{"kind":"MANEUVER","context":{"value":9007199254740992}}`},
		{name: "fraction", payload: `{"kind":"MANEUVER","context":{"value":1.5}}`},
		{name: "unpaired surrogate", payload: `{"kind":"MANEUVER","source_ref":"\ud800"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := contracts.Command{
				ID:               model.CommandID("invalid-" + test.name),
				SchemaVersion:    CommandSchemaV1,
				Type:             CommandStartAction,
				MatchID:          "match-1",
				ActorPlayerID:    "player-1",
				ExpectedRevision: &revision,
				Payload:          json.RawMessage(test.payload),
			}
			if _, _, _, err := normalizeCommandRequest("principal-1", command); CodeOf(err) != CodeInvalidCommand {
				t.Fatalf("expected invalid command, got %v", err)
			}
		})
	}
}
