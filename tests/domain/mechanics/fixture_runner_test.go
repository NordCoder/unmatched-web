package mechanics_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/rules"
)

type fixtureSuite struct {
	SchemaVersion string        `json:"schema_version"`
	Cases         []fixtureCase `json:"cases"`
}
type fixtureCase struct {
	Name                 string                     `json:"name"`
	Definition           effects.Definition         `json:"definition"`
	State                model.GameState            `json:"state"`
	Procedure            model.ProcedureRef         `json:"procedure"`
	ExpectedFirstStatus  contracts.ResolutionStatus `json:"expected_first_status"`
	ExpectedFirstEvents  []string                   `json:"expected_first_events"`
	ChoiceOptionIndex    *int                       `json:"choice_option_index,omitempty"`
	ExpectedResumeEvents []string                   `json:"expected_resume_events,omitempty"`
	ExpectedFinalStatus  contracts.ResolutionStatus `json:"expected_final_status"`
	ExpectedFighterState *struct {
		Fighter model.FighterID `json:"fighter"`
		Key     string          `json:"key"`
		Value   any             `json:"value"`
	} `json:"expected_fighter_state,omitempty"`
}
type pathResult struct {
	first  contracts.ResolutionOutcome
	final  contracts.ResolutionOutcome
	events []contracts.DomainEvent
	state  model.GameState
}

func TestI1MachineReadableFixtures(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	fixtureBytes, err := os.ReadFile(filepath.Join(filepath.Dir(file), "..", "..", "fixtures", "mechanics", "i1.json"))
	if err != nil {
		t.Fatal(err)
	}
	var suite fixtureSuite
	if err = json.Unmarshal(fixtureBytes, &suite); err != nil {
		t.Fatal(err)
	}
	if suite.SchemaVersion != "unmatched.mechanics.fixture/v2" {
		t.Fatalf("unexpected schema %q", suite.SchemaVersion)
	}
	if len(suite.Cases) < 4 {
		t.Fatal("incomplete I1 suite")
	}
	for _, tc := range suite.Cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			engine, err := rules.New([]effects.Definition{tc.Definition})
			if err != nil {
				t.Fatalf("load definition: %v", err)
			}
			initialBytes := canonical(t, tc.State)
			direct := runFixturePath(t, engine, tc, initialBytes, false)
			serialized := runFixturePath(t, engine, tc, initialBytes, true)
			assertOutcome(t, direct.first, tc.ExpectedFirstStatus, tc.ExpectedFirstEvents)
			assertOutcome(t, serialized.first, tc.ExpectedFirstStatus, tc.ExpectedFirstEvents)
			assertOutcome(t, serialized.final, tc.ExpectedFinalStatus, tc.ExpectedResumeEvents)
			if !reflect.DeepEqual(canonical(t, direct.state), canonical(t, serialized.state)) {
				t.Fatal("direct and serialized pause/resume states diverged")
			}
			replayed := restoreState(t, initialBytes)
			for _, ev := range serialized.events {
				replayed, err = engine.ApplyEvent(replayed, ev)
				if err != nil {
					t.Fatal(err)
				}
			}
			if !reflect.DeepEqual(canonical(t, replayed), canonical(t, serialized.state)) {
				t.Fatal("full event replay diverged")
			}
			repeat := runFixturePath(t, engine, tc, initialBytes, true)
			if !reflect.DeepEqual(canonical(t, repeat.state), canonical(t, serialized.state)) {
				t.Fatal("repeated execution diverged")
			}
			if !reflect.DeepEqual(canonical(t, tc.State), initialBytes) {
				t.Fatal("fixture state mutated")
			}
			if tc.ExpectedFighterState != nil {
				got := serialized.state.Fighters[tc.ExpectedFighterState.Fighter].State[tc.ExpectedFighterState.Key]
				if !reflect.DeepEqual(got, tc.ExpectedFighterState.Value) {
					t.Fatalf("state %v != %v", got, tc.ExpectedFighterState.Value)
				}
			}
		})
	}
}

func runFixturePath(t *testing.T, engine *rules.Engine, tc fixtureCase, initial []byte, serialized bool) pathResult {
	t.Helper()
	state := restoreState(t, initial)
	input := contracts.ResolutionInput{CommandID: model.CommandID("fixture:" + tc.Name), Procedure: tc.Procedure}
	first, err := engine.Resolve(state, input)
	if err != nil {
		t.Fatal(err)
	}
	events := append([]contracts.DomainEvent(nil), first.Events...)
	state = reduceAll(t, engine, state, first.Events)
	final := first
	if tc.ChoiceOptionIndex != nil {
		if first.PendingInteraction == nil {
			t.Fatal("fixture expected interaction")
		}
		var domain struct {
			Options []struct {
				Handle string `json:"handle"`
			} `json:"options"`
		}
		if err = json.Unmarshal(first.PendingInteraction.LegalDomain, &domain); err != nil {
			t.Fatal(err)
		}
		idx := *tc.ChoiceOptionIndex
		if idx < 0 || idx >= len(domain.Options) {
			t.Fatal("invalid fixture option index")
		}
		proc := first.PendingInteraction.ResumeProcedure
		if serialized {
			b := canonical(t, proc)
			if err = json.Unmarshal(b, &proc); err != nil {
				t.Fatal(err)
			}
		}
		final, err = engine.Resolve(state, contracts.ResolutionInput{CommandID: model.CommandID("fixture-resume:" + tc.Name), Procedure: proc, InteractionID: first.PendingInteraction.ID, Choice: raw(domain.Options[idx].Handle)})
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, final.Events...)
		state = reduceAll(t, engine, state, final.Events)
	}
	return pathResult{first: first, final: final, events: events, state: state}
}
func canonical(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func restoreState(t *testing.T, b []byte) model.GameState {
	t.Helper()
	var s model.GameState
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatal(err)
	}
	return s
}
func assertOutcome(t *testing.T, out contracts.ResolutionOutcome, status contracts.ResolutionStatus, types []string) {
	t.Helper()
	if out.Status != status {
		t.Fatalf("status %s != %s", out.Status, status)
	}
	if types == nil {
		return
	}
	got := make([]string, len(out.Events))
	for i, e := range out.Events {
		got[i] = e.Type
	}
	if !reflect.DeepEqual(got, types) {
		t.Fatalf("events %v != %v", got, types)
	}
}
