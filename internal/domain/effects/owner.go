package effects

import (
	"encoding/json"
	"fmt"

	"github.com/NordCoder/unmatched-web/internal/domain/query"
)

const TrustedActorBinding = "actor_player_id"

type OwnerSelector string

const OwnerActor OwnerSelector = "actor"

// ActorOwner returns the only supported Wave 1 interaction owner selector.
// It is represented internally as a closed reference to the trusted host
// binding, but serialized definitions expose only owner_selector: "actor".
func ActorOwner() query.Expr {
	return query.Expr{
		Kind:       query.Reference,
		Source:     query.Captured,
		Path:       []string{TrustedActorBinding},
		ValueType:  query.TypePlayer,
		Visibility: query.Public,
	}
}

func IsActorOwner(expr query.Expr) bool {
	return expr.Kind == query.Reference &&
		expr.Source == query.Captured &&
		len(expr.Path) == 1 && expr.Path[0] == TrustedActorBinding &&
		expr.ValueType == query.TypePlayer &&
		(expr.Visibility == "" || expr.Visibility == query.Public) &&
		expr.Value == nil && len(expr.Args) == 0 && len(expr.Fields) == 0
}

func ValidateOwnerSelector(expr query.Expr) error {
	if !IsActorOwner(expr) {
		return fmt.Errorf("Wave 1 choice owner must use owner_selector actor")
	}
	return nil
}

// MarshalJSON keeps the legacy Go field private to implementation code and
// publishes the reduced Wave 1 definition contract.
func (c Choice) MarshalJSON() ([]byte, error) {
	type wireChoice struct {
		Kind          string            `json:"kind"`
		Binding       string            `json:"binding"`
		Visibility    query.Visibility  `json:"visibility"`
		OwnerSelector OwnerSelector     `json:"owner_selector,omitempty"`
		Owner         *query.Expr       `json:"owner,omitempty"`
		Domain        query.Expr        `json:"domain"`
		Prompt        any               `json:"prompt"`
		EmptyDomain   EmptyDomainPolicy `json:"empty_domain"`
		ValueType     query.Type        `json:"value_type,omitempty"`
		Default       *query.Expr       `json:"default,omitempty"`
		Multi         bool              `json:"multi,omitempty"`
	}
	w := wireChoice{
		Kind: c.Kind, Binding: c.Binding, Visibility: c.Visibility,
		Domain: c.Domain, Prompt: c.Prompt, EmptyDomain: c.EmptyDomain,
		ValueType: c.ValueType, Default: c.Default, Multi: c.Multi,
	}
	if IsActorOwner(c.Owner) {
		w.OwnerSelector = OwnerActor
	} else {
		legacy := c.Owner
		w.Owner = &legacy
	}
	return json.Marshal(w)
}

// UnmarshalJSON accepts the legacy owner field only so rules.New can reject it
// explicitly. No legacy expression is converted into actor ownership.
func (c *Choice) UnmarshalJSON(data []byte) error {
	type wireChoice struct {
		Kind          string            `json:"kind"`
		Binding       string            `json:"binding"`
		Visibility    query.Visibility  `json:"visibility"`
		OwnerSelector OwnerSelector     `json:"owner_selector,omitempty"`
		Owner         *query.Expr       `json:"owner,omitempty"`
		Domain        query.Expr        `json:"domain"`
		Prompt        any               `json:"prompt"`
		EmptyDomain   EmptyDomainPolicy `json:"empty_domain"`
		ValueType     query.Type        `json:"value_type,omitempty"`
		Default       *query.Expr       `json:"default,omitempty"`
		Multi         bool              `json:"multi,omitempty"`
	}
	var w wireChoice
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	if w.OwnerSelector != "" && w.Owner != nil {
		return fmt.Errorf("choice cannot contain both owner_selector and legacy owner")
	}
	*c = Choice{
		Kind: w.Kind, Binding: w.Binding, Visibility: w.Visibility,
		Domain: w.Domain, Prompt: w.Prompt, EmptyDomain: w.EmptyDomain,
		ValueType: w.ValueType, Default: w.Default, Multi: w.Multi,
	}
	switch w.OwnerSelector {
	case OwnerActor:
		c.Owner = ActorOwner()
	case "":
		if w.Owner != nil {
			c.Owner = *w.Owner
		}
	default:
		return fmt.Errorf("unsupported owner_selector %q", w.OwnerSelector)
	}
	return nil
}
