package phase4bcards

import (
	"bytes"
	"testing"
)

func TestRuntimeManifestStopsAfterCardsSection(t *testing.T) {
	raw := []byte("schema_version: 1\ncards:\n  - {id: card-one}\nsources:\n  - {id: source-record}\nvalidation:\n  quantity: pass\n")
	got := runtimeManifest(raw)
	if !bytes.Contains(got, []byte("card-one")) {
		t.Fatal("cards section was removed")
	}
	if bytes.Contains(got, []byte("source-record")) || bytes.Contains(got, []byte("validation:")) {
		t.Fatalf("non-card sections leaked into runtime manifest: %q", got)
	}
}
