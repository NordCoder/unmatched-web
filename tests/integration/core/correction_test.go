package core_test

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/application"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	coreruntime "github.com/NordCoder/unmatched-web/internal/runtime"
	"github.com/NordCoder/unmatched-web/internal/transport"
)

func TestProjectedInteractionIsAllowListedBeforeTransportEncoding(t *testing.T) {
	rules := &secretPendingRules{}
	host, _ := newHost(t, coreruntime.NewSequenceIDProvider("projection"), rules)
	created, _ := createAndJoin(t, host)
	two := uint64(2)
	result, err := host.Execute(context.Background(), "principal-a", command(
		"action-secret", application.CommandStartAction, created.MatchID, created.PlayerID, &two,
		application.StartActionPayload{
			Kind:    application.ActionScheme,
			Context: map[string]json.RawMessage{"internal_secret": json.RawMessage(`"do-not-deliver"`)},
		},
	))
	if err != nil {
		t.Fatalf("start action: %v", err)
	}
	encoded, err := transport.EncodeProjection(result.Projection)
	if err != nil {
		t.Fatalf("encode projection: %v", err)
	}
	for _, forbidden := range [][]byte{
		[]byte("resume_procedure"), []byte("procedure_instance_id"),
		[]byte("bindings"), []byte("internal_secret"), []byte("do-not-deliver"),
	} {
		if bytes.Contains(encoded, forbidden) {
			t.Fatalf("encoded projection leaked internal resolver field %q: %s", forbidden, encoded)
		}
	}
	var envelope struct {
		Projection struct {
			Pending map[string]json.RawMessage `json:"pending_interaction"`
		} `json:"projection"`
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	expected := []string{"interaction_instance_id", "owner_player_id", "kind", "visibility", "prompt", "legal_domain"}
	if len(envelope.Projection.Pending) != len(expected) {
		t.Fatalf("unexpected projected interaction shape: %v", envelope.Projection.Pending)
	}
	for _, key := range expected {
		if _, ok := envelope.Projection.Pending[key]; !ok {
			t.Fatalf("projected interaction omitted %q", key)
		}
	}
}

func TestCommandIdentityIsScopedByPrincipal(t *testing.T) {
	host, _ := newHost(t, coreruntime.NewSequenceIDProvider("scope"), &completedRules{})
	request := command(
		"shared-command-id", application.CommandCreateMatch, "", "", nil,
		application.CreateMatchPayload{DefinitionKey: "synthetic@v1", FighterDefinition: "mirror"},
	)
	first, err := host.Execute(context.Background(), "principal-a", request)
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	second, err := host.Execute(context.Background(), "principal-b", request)
	if err != nil {
		t.Fatalf("second principal create: %v", err)
	}
	if first.MatchID == second.MatchID || first.PlayerID == second.PlayerID {
		t.Fatalf("independent principal command identities reused runtime IDs: first=%+v second=%+v", first, second)
	}
}

func TestConcurrentDuplicateDoesNotRepeatRulesOrIDAllocation(t *testing.T) {
	ids := &countingIDs{delegate: coreruntime.NewSequenceIDProvider("concurrent")}
	rules := newBlockingPendingRules()
	host, store := newHost(t, ids, rules)
	created, _ := createAndJoin(t, host)
	baselineIDs := ids.Count()
	two := uint64(2)
	request := command(
		"same-action", application.CommandStartAction, created.MatchID, created.PlayerID, &two,
		application.StartActionPayload{Kind: application.ActionScheme},
	)

	const callers = 8
	results := make([]application.CommandResult, callers)
	errs := make([]error, callers)
	var wg sync.WaitGroup
	wg.Add(callers)
	for index := 0; index < callers; index++ {
		go func(index int) {
			defer wg.Done()
			results[index], errs[index] = host.Execute(context.Background(), "principal-a", request)
		}(index)
	}
	<-rules.entered
	close(rules.release)
	wg.Wait()

	for index, err := range errs {
		if err != nil {
			t.Fatalf("caller %d failed: %v", index, err)
		}
		if !reflect.DeepEqual(results[0], results[index]) {
			t.Fatalf("caller %d received a different duplicate result", index)
		}
	}
	if got := rules.ResolveCalls(); got != 1 {
		t.Fatalf("RulesEngine invoked %d times, want 1", got)
	}
	// One procedure, one action, and two event IDs are allocated for the accepted
	// action. Waiting duplicates must allocate nothing.
	if delta := ids.Count() - baselineIDs; delta != 4 {
		t.Fatalf("duplicate callers allocated %d IDs, want 4 for one execution", delta)
	}
	if events := store.Events(created.MatchID); len(events) != 4 {
		t.Fatalf("duplicate callers changed event count to %d, want 4", len(events))
	}
}

type secretPendingRules struct{}

func (r *secretPendingRules) Resolve(_ model.GameState, input contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	var owner model.PlayerID
	if err := json.Unmarshal(input.Procedure.Bindings["actor_player_id"], &owner); err != nil {
		return contracts.ResolutionOutcome{}, err
	}
	return contracts.ResolutionOutcome{
		Status: contracts.ResolutionPending,
		PendingInteraction: &model.PendingInteraction{
			ID: "secret-interaction", OwnerPlayerID: owner, Kind: "SECRET_CHOICE", Visibility: "OWNER_ONLY",
			Prompt: json.RawMessage(`{"message":"safe prompt"}`), LegalDomain: json.RawMessage(`{"options":[1,2]}`),
			ResumeProcedure: input.Procedure,
		},
	}, nil
}
func (r *secretPendingRules) LegalActions(model.GameState, model.PlayerID) ([]json.RawMessage, error) {
	return nil, nil
}
func (r *secretPendingRules) Project(_ model.GameState, _ model.PlayerID) (json.RawMessage, error) {
	return json.RawMessage(`{"safe":true}`), nil
}

type completedRules struct{}

func (r *completedRules) Resolve(model.GameState, contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	return contracts.ResolutionOutcome{Status: contracts.ResolutionCompleted}, nil
}
func (r *completedRules) LegalActions(model.GameState, model.PlayerID) ([]json.RawMessage, error) {
	return nil, nil
}
func (r *completedRules) Project(_ model.GameState, _ model.PlayerID) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}

type countingIDs struct {
	delegate coreruntime.IDProvider
	count    atomic.Int64
}

func (p *countingIDs) Next(kind string) string { p.count.Add(1); return p.delegate.Next(kind) }
func (p *countingIDs) Count() int64            { return p.count.Load() }

type blockingPendingRules struct {
	entered chan struct{}
	release chan struct{}
	once    sync.Once
	calls   atomic.Int64
}

func newBlockingPendingRules() *blockingPendingRules {
	return &blockingPendingRules{entered: make(chan struct{}), release: make(chan struct{})}
}
func (r *blockingPendingRules) Resolve(_ model.GameState, input contracts.ResolutionInput) (contracts.ResolutionOutcome, error) {
	r.calls.Add(1)
	r.once.Do(func() { close(r.entered) })
	<-r.release
	var owner model.PlayerID
	if err := json.Unmarshal(input.Procedure.Bindings["actor_player_id"], &owner); err != nil {
		return contracts.ResolutionOutcome{}, err
	}
	return contracts.ResolutionOutcome{
		Status: contracts.ResolutionPending,
		PendingInteraction: &model.PendingInteraction{
			ID: "concurrent-interaction", OwnerPlayerID: owner, Kind: "CHOICE", Visibility: "OWNER_ONLY",
			Prompt: json.RawMessage(`{}`), LegalDomain: json.RawMessage(`{}`), ResumeProcedure: input.Procedure,
		},
	}, nil
}
func (r *blockingPendingRules) ResolveCalls() int64 { return r.calls.Load() }
func (r *blockingPendingRules) LegalActions(model.GameState, model.PlayerID) ([]json.RawMessage, error) {
	return nil, nil
}
func (r *blockingPendingRules) Project(_ model.GameState, _ model.PlayerID) (json.RawMessage, error) {
	return json.RawMessage(`{}`), nil
}
