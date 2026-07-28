package rules

import (
	"encoding/json"
	"fmt"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
)

func validateWave1Safety(d effects.Definition) error {
	for name, spec := range d.CapturedBindings {
		if err := query.ValidateClosedSpec(spec); err != nil {
			return fmt.Errorf("captured binding %s: %w", name, err)
		}
		if input, ok := d.InputBindings[name]; ok && !query.ExactSpec(spec, input) {
			return fmt.Errorf("binding %s has conflicting captured/input specs", name)
		}
	}
	for name, spec := range d.InputBindings {
		if err := query.ValidateClosedSpec(spec); err != nil {
			return fmt.Errorf("input binding %s: %w", name, err)
		}
	}
	for _, stage := range d.Stages {
		if stage.Choice == nil {
			continue
		}
		if stage.Choice.Multi {
			return fmt.Errorf("choice %s multi-select is not supported in Wave 1", stage.Choice.Binding)
		}
		if stage.Choice.EmptyDomain == effects.EmptyBindEmpty {
			return fmt.Errorf("choice %s bind_empty is not supported in Wave 1", stage.Choice.Binding)
		}
	}
	return nil
}

func capturedTransportSpecs(d effects.Definition) map[string]query.ValueSpec {
	specs := make(map[string]query.ValueSpec, len(d.CapturedBindings)+len(d.InputBindings))
	for name, spec := range d.CapturedBindings {
		specs[name] = spec
	}
	for name, spec := range d.InputBindings {
		if _, captured := specs[name]; captured {
			continue
		}
		spec.Optional = true
		specs[name] = spec
	}
	return specs
}

func rejectWithDiagnostic(code string, err error) contracts.ResolutionOutcome {
	return contracts.ResolutionOutcome{Status: contracts.ResolutionRejected, RejectionCode: code, Diagnostics: map[string]string{"bindings": err.Error()}}
}

func canonicalRawEqual(a, b json.RawMessage) bool {
	var av, bv any
	if json.Unmarshal(a, &av) != nil || json.Unmarshal(b, &bv) != nil {
		return false
	}
	ar, _ := json.Marshal(av)
	br, _ := json.Marshal(bv)
	return string(ar) == string(br)
}
