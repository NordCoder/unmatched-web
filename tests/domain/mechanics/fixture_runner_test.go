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
	if suite.SchemaVersion != "unmatched.mechanics.fixture/v1" {
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
			input := contracts.ResolutionInput{CommandID: model.CommandID("fixture:" + tc.Name), Procedure: tc.Procedure}
			first, err := engine.Resolve(tc.State, input)
			if err != nil {
				t.Fatal(err)
			}
			assertOutcome(t, first, tc.ExpectedFirstStatus, tc.ExpectedFirstEvents)
			encoded, _ := json.Marshal(first)
			var restored contracts.ResolutionOutcome
			if err = json.Unmarshal(encoded, &restored); err != nil {
				t.Fatal(err)
			}
			persisted := tc.State
			for _, event := range restored.Events {
				persisted, err = engine.ApplyEvent(persisted, event)
				if err != nil {
					t.Fatal(err)
				}
			}
			final := restored
			if tc.ChoiceOptionIndex != nil {
				if restored.PendingInteraction == nil {
					t.Fatal("fixture expected interaction")
				}
				var domain struct {
					Options []struct {
						Handle string `json:"handle"`
					} `json:"options"`
				}
				if err = json.Unmarshal(restored.PendingInteraction.LegalDomain, &domain); err != nil {
					t.Fatal(err)
				}
				idx := *tc.ChoiceOptionIndex
				if idx < 0 || idx >= len(domain.Options) {
					t.Fatal("invalid fixture option index")
				}
				procedureBytes, _ := json.Marshal(restored.PendingInteraction.ResumeProcedure)
				var resumed model.ProcedureRef
				if err = json.Unmarshal(procedureBytes, &resumed); err != nil {
					t.Fatal(err)
				}
				final, err = engine.Resolve(persisted, contracts.ResolutionInput{CommandID: model.CommandID("fixture-resume:" + tc.Name), Procedure: resumed, InteractionID: restored.PendingInteraction.ID, Choice: raw(domain.Options[idx].Handle)})
				if err != nil {
					t.Fatal(err)
				}
				assertOutcome(t, final, tc.ExpectedFinalStatus, tc.ExpectedResumeEvents)
			} else if final.Status != tc.ExpectedFinalStatus {
				t.Fatalf("final status %s", final.Status)
			}
			if tc.ExpectedFighterState != nil {
				got := persisted.Fighters[tc.ExpectedFighterState.Fighter].State[tc.ExpectedFighterState.Key]
				if !reflect.DeepEqual(got, tc.ExpectedFighterState.Value) {
					t.Fatalf("state %v != %v", got, tc.ExpectedFighterState.Value)
				}
			}
			// Re-running from identical bytes must be deterministic.
			again, err := engine.Resolve(tc.State, input)
			if err != nil {
				t.Fatal(err)
			}
			a, _ := json.Marshal(first)
			b, _ := json.Marshal(again)
			if !reflect.DeepEqual(a, b) {
				t.Fatal("fixture replay divergence")
			}
		})
	}
}
func assertOutcome(t *testing.T, out contracts.ResolutionOutcome, status contracts.ResolutionStatus, types []string) {
	t.Helper()
	if out.Status != status {
		t.Fatalf("status %s != %s", out.Status, status)
	}
	got := make([]string, len(out.Events))
	for i, e := range out.Events {
		got[i] = e.Type
	}
	if !reflect.DeepEqual(got, types) {
		t.Fatalf("events %v != %v", got, types)
	}
}
