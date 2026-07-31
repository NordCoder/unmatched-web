// Package server exposes the minimal local HTTP boundary for a two-tab
// playable match. It intentionally uses in-memory match ownership while the
// approved PostgreSQL Core adapter remains unchanged and available elsewhere.
package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/NordCoder/unmatched-web/internal/playableslice/content"
	"github.com/NordCoder/unmatched-web/internal/playableslice/game"
)

type Manager struct {
	mu       sync.Mutex
	registry *content.Registry
	matches  map[string]*game.Match
	sessions map[string]session
}

type session struct {
	MatchID  string
	PlayerID string
}

type apiResponse struct {
	MatchID  string     `json:"match_id,omitempty"`
	Token    string     `json:"token,omitempty"`
	PlayerID string     `json:"player_id,omitempty"`
	View     *game.View `json:"view,omitempty"`
	Error    *apiError  `json:"error,omitempty"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewManager(registry *content.Registry) *Manager {
	return &Manager{registry: registry, matches: map[string]*game.Match{}, sessions: map[string]session{}}
}

func NewHandler(registry *content.Registry, static fs.FS) http.Handler {
	manager := NewManager(registry)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("POST /api/matches", manager.createMatch)
	mux.HandleFunc("POST /api/matches/{matchID}/join", manager.joinMatch)
	mux.HandleFunc("GET /api/matches/{matchID}", manager.getMatch)
	mux.HandleFunc("POST /api/matches/{matchID}/commands", manager.command)
	if static != nil {
		fileServer := http.FileServer(http.FS(static))
		mux.Handle("/", fileServer)
	}
	return securityHeaders(requestLog(mux))
}

func (m *Manager) createMatch(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	matchID, err := randomID("match", 8)
	if err != nil {
		writeInternal(w, err)
		return
	}
	match, playerID, err := game.New(matchID, m.registry, "robin-hood")
	if err != nil {
		writeInternal(w, err)
		return
	}
	token, err := randomID("session", 16)
	if err != nil {
		writeInternal(w, err)
		return
	}
	m.matches[matchID] = match
	m.sessions[token] = session{MatchID: matchID, PlayerID: playerID}
	view, _ := match.Project(playerID)
	writeJSON(w, http.StatusCreated, apiResponse{MatchID: matchID, Token: token, PlayerID: playerID, View: &view})
}

func (m *Manager) joinMatch(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	matchID := r.PathValue("matchID")
	match := m.matches[matchID]
	if match == nil {
		writeAPIError(w, http.StatusNotFound, "match_not_found", "match does not exist")
		return
	}
	playerID, err := match.Join("bigfoot")
	if err != nil {
		writeRuleError(w, err)
		return
	}
	token, err := randomID("session", 16)
	if err != nil {
		writeInternal(w, err)
		return
	}
	m.sessions[token] = session{MatchID: matchID, PlayerID: playerID}
	view, _ := match.Project(playerID)
	writeJSON(w, http.StatusOK, apiResponse{MatchID: matchID, Token: token, PlayerID: playerID, View: &view})
}

func (m *Manager) getMatch(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	session, match, ok := m.authorize(r, r.PathValue("matchID"))
	if !ok {
		writeAPIError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid match token")
		return
	}
	view, err := match.Project(session.PlayerID)
	if err != nil {
		writeRuleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{MatchID: match.ID, PlayerID: session.PlayerID, View: &view})
}

func (m *Manager) command(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	session, match, ok := m.authorize(r, r.PathValue("matchID"))
	if !ok {
		writeAPIError(w, http.StatusUnauthorized, "unauthorized", "missing or invalid match token")
		return
	}
	var command game.Command
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&command); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := match.Apply(session.PlayerID, command); err != nil {
		writeRuleError(w, err)
		return
	}
	view, err := match.Project(session.PlayerID)
	if err != nil {
		writeRuleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{MatchID: match.ID, PlayerID: session.PlayerID, View: &view})
}

func (m *Manager) authorize(r *http.Request, matchID string) (session, *game.Match, bool) {
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	s, ok := m.sessions[token]
	if !ok || s.MatchID != matchID {
		return session{}, nil, false
	}
	match := m.matches[matchID]
	if match == nil {
		return session{}, nil, false
	}
	return s, match, true
}

func randomID(prefix string, bytesCount int) (string, error) {
	buffer := make([]byte, bytesCount)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return prefix + "-" + hex.EncodeToString(buffer), nil
}

func writeRuleError(w http.ResponseWriter, err error) {
	var rule *game.RuleError
	if errors.As(err, &rule) {
		status := http.StatusConflict
		if rule.Code == "unauthorized" || rule.Code == "unauthorized_choice" {
			status = http.StatusForbidden
		}
		writeAPIError(w, status, rule.Code, rule.Message)
		return
	}
	writeInternal(w, err)
}
func writeInternal(w http.ResponseWriter, err error) {
	slog.Error("playable slice request failed", "error", err)
	writeAPIError(w, http.StatusInternalServerError, "internal_error", "internal server error")
}
func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSONStatus(w, status, apiResponse{Error: &apiError{Code: code, Message: message}})
}
func writeJSON(w http.ResponseWriter, status int, value any) { writeJSONStatus(w, status, value) }
func writeJSONStatus(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; connect-src 'self'; img-src 'self' data:")
		next.ServeHTTP(w, r)
	})
}
func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("request", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func RunAddress(port string) string {
	if strings.TrimSpace(port) == "" {
		port = "8080"
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return fmt.Sprintf(":%s", port)
}
