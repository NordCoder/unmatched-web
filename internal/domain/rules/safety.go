package rules

import (
	"encoding/json"
	"fmt"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/effects"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
	"github.com/NordCoder/unmatched-web/internal/domain/query"
)

const actorPlayerIDBinding = effects.TrustedActorBinding

var trustedHostBindingSpecs = map[string]query.ValueSpec{
	actorPlayerIDBinding: {Type: query.TypePlayer, Visibility: query.Public, Optional: true},
}

func validateWave1Safety(d effects.Definition) error {
	for name, spec := range d.CapturedBindings {
		if _, reserved := trustedHostBindingSpecs[name]; reserved {
			return fmt.Errorf("captured binding %s is reserved for the trusted host", name)
		}
		if err := query.ValidateClosedSpec(spec); err != nil {
			return fmt.Errorf("captured binding %s: %w", name, err)
		}
		if input, ok := d.InputBindings[name]; ok && !query.ExactSpec(spec, input) {
			return fmt.Errorf("binding %s has conflicting captured/input specs", name)
		}
	}
	for name, spec := range d.InputBindings {
		if _, reserved := trustedHostBindingSpecs[name]; reserved {
			return fmt.Errorf("input binding %s is reserved for the trusted host", name)
		}
		if err := query.ValidateClosedSpec(spec); err != nil {
			return fmt.Errorf("input binding %s: %w", name, err)
		}
	}
	for _, stage := range d.Stages {
		if stage.Choice == nil {
			continue
		}
		if err := effects.ValidateOwnerSelector(stage.Choice.Owner); err != nil {
			return fmt.Errorf("choice %s: %w", stage.Choice.Binding, err)
		}
		if stage.Choice.Multi {
			return fmt.Errorf("choice %s multi-select is not supported in Wave 1", stage.Choice.Binding)
		}
		if stage.Choice.EmptyDomain == effects.EmptyBindEmpty {
			return fmt.Errorf("choice %s bind_empty is not supported in Wave 1", stage.Choice.Binding)
		}
	}
	return validatePublicControlFlow(d)
}

func definitionForValidation(d effects.Definition) effects.Definition {
	if !requiresActorOwner(d) {
		return d
	}
	d.Stages = append([]effects.Stage(nil), d.Stages...)
	for i := range d.Stages {
		if d.Stages[i].Choice == nil {
			continue
		}
		choice := *d.Stages[i].Choice
		if effects.IsActorOwner(choice.Owner) {
			// effects.Validate still has a legacy owner-expression slot. Feed it a
			// validation-only public player literal so the trusted actor binding is
			// not added to the general definition expression environment.
			choice.Owner = query.Expr{Kind: query.Literal, Value: "__wave1_actor__", ValueType: query.TypePlayer, Visibility: query.Public}
		}
		d.Stages[i].Choice = &choice
	}
	return d
}

func requiresActorOwner(d effects.Definition) bool {
	for _, stage := range d.Stages {
		if stage.Choice != nil && effects.IsActorOwner(stage.Choice.Owner) {
			return true
		}
	}
	return false
}

func validatePublicControlFlow(d effects.Definition) error {
	registry := operations.Default()
	env := query.TypeEnv{
		Captured: copyBindingSpecs(d.CapturedBindings),
		Results:  map[string]query.ValueSpec{},
		Choices:  map[string]query.ValueSpec{},
		Input:    copyBindingSpecs(d.InputBindings),
	}
	for _, stage := range d.Stages {
		if stage.Condition != nil {
			spec, err := query.Infer(*stage.Condition, env)
			if err != nil {
				return fmt.Errorf("stage %s condition: %w", stage.ID, err)
			}
			if spec.Visibility != query.Public {
				return fmt.Errorf("stage %s condition must be public", stage.ID)
			}
		}
		for _, prerequisite := range stage.Prerequisites {
			spec, err := query.Infer(prerequisite, env)
			if err != nil {
				return fmt.Errorf("stage %s prerequisite: %w", stage.ID, err)
			}
			if spec.Visibility != query.Public {
				return fmt.Errorf("stage %s prerequisite must be public", stage.ID)
			}
		}
		for _, operation := range stage.Costs {
			if operation.ResultBinding == "" {
				continue
			}
			spec, err := registry.Validate(operation, env)
			if err != nil {
				return err
			}
			env.Results[operation.ResultBinding] = spec
		}
		if stage.Choice != nil {
			spec, err := effects.ChoiceSpec(d, stage.Choice.Binding, registry)
			if err != nil {
				return err
			}
			if stage.Choice.EmptyDomain == effects.EmptyReject || stage.Choice.EmptyDomain == effects.EmptyBindDefault {
				env.Choices[stage.Choice.Binding] = spec
			}
		}
		for _, operation := range stage.Operations {
			if operation.ResultBinding == "" {
				continue
			}
			spec, err := registry.Validate(operation, env)
			if err != nil {
				return err
			}
			env.Results[operation.ResultBinding] = spec
		}
	}
	return nil
}

func copyBindingSpecs(source map[string]query.ValueSpec) map[string]query.ValueSpec {
	result := make(map[string]query.ValueSpec, len(source)+1)
	for name, spec := range source {
		result[name] = spec
	}
	return result
}

func capturedTransportSpecs(d effects.Definition) map[string]query.ValueSpec {
	specs := make(map[string]query.ValueSpec, len(d.CapturedBindings)+len(d.InputBindings)+len(trustedHostBindingSpecs))
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
	for name, spec := range trustedHostBindingSpecs {
		if name == actorPlayerIDBinding && requiresActorOwner(d) {
			spec.Optional = false
		}
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
