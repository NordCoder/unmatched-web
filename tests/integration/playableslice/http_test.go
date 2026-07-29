package playableslice_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
	"github.com/NordCoder/unmatched-web/internal/playableslice/game"
	sliceserver "github.com/NordCoder/unmatched-web/internal/playableslice/server"
)

type response struct {
	MatchID  string                          `json:"match_id"`
	Token    string                          `json:"token"`
	PlayerID string                          `json:"player_id"`
	View     *game.View                      `json:"view"`
	Error    *struct{ Code, Message string } `json:"error"`
}

func TestTwoBrowserSessionBoundary(t *testing.T) {
	registry, err := content.Load()
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(sliceserver.NewHandler(registry, nil))
	defer server.Close()
	first := do(t, http.MethodPost, server.URL+"/api/matches", "", map[string]any{})
	if first.Token == "" || first.MatchID == "" {
		t.Fatalf("invalid create response: %+v", first)
	}
	second := do(t, http.MethodPost, server.URL+"/api/matches/"+first.MatchID+"/join", "", map[string]any{})
	if second.Token == "" || second.PlayerID == first.PlayerID {
		t.Fatalf("invalid join response: %+v", second)
	}
	firstState := do(t, http.MethodGet, server.URL+"/api/matches/"+first.MatchID, first.Token, nil)
	secondState := do(t, http.MethodGet, server.URL+"/api/matches/"+first.MatchID, second.Token, nil)
	if firstState.View.Phase != game.PhaseActive || secondState.View.Phase != game.PhaseActive {
		t.Fatal("match did not become active")
	}
	assertPrivateHand(t, firstState.View, first.PlayerID)
	assertPrivateHand(t, secondState.View, second.PlayerID)

	unauthorized, code := doStatus(t, http.MethodGet, server.URL+"/api/matches/"+first.MatchID, "wrong-token", nil)
	if code != http.StatusUnauthorized || unauthorized.Error == nil || unauthorized.Error.Code != "unauthorized" {
		t.Fatalf("unauthorized status=%d response=%+v", code, unauthorized)
	}
	revision := firstState.View.Revision
	stale, code := doStatus(t, http.MethodPost, server.URL+"/api/matches/"+first.MatchID+"/commands", first.Token, map[string]any{"type": "maneuver", "expected_revision": revision + 10})
	if code != http.StatusConflict || stale.Error == nil || stale.Error.Code != "stale_revision" {
		t.Fatalf("stale status=%d response=%+v", code, stale)
	}
	readback := do(t, http.MethodGet, server.URL+"/api/matches/"+first.MatchID, first.Token, nil)
	if readback.View.Revision != revision {
		t.Fatal("rejected HTTP command mutated match")
	}
}

func assertPrivateHand(t *testing.T, view *game.View, viewer string) {
	t.Helper()
	for _, player := range view.Players {
		if player.ID == viewer {
			if len(player.Hand) != 5 {
				t.Fatalf("viewer hand=%d", len(player.Hand))
			}
		} else if len(player.Hand) != 0 {
			t.Fatal("opponent hand leaked")
		}
	}
}
func do(t *testing.T, method, url, token string, body any) response {
	t.Helper()
	value, status := doStatus(t, method, url, token, body)
	if status < 200 || status >= 300 {
		t.Fatalf("%s %s status=%d error=%+v", method, url, status, value.Error)
	}
	return value
}
func doStatus(t *testing.T, method, url, token string, body any) (response, int) {
	t.Helper()
	var raw []byte
	if body != nil {
		raw, _ = json.Marshal(body)
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var value response
	if err := json.NewDecoder(res.Body).Decode(&value); err != nil {
		t.Fatal(err)
	}
	return value, res.StatusCode
}
