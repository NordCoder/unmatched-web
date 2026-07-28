// Package history builds deterministic, rebuildable provenance indexes from actual rules events.
package history

import (
	"encoding/json"
	"sort"

	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"github.com/NordCoder/unmatched-web/internal/domain/operations"
)

type Entry struct {
	EventID      model.EventID `json:"event_id"`
	Sequence     uint64        `json:"sequence"`
	Type         string        `json:"type"`
	SourceRef    string        `json:"source_ref"`
	OperationID  string        `json:"operation_instance_id"`
	CauseKind    string        `json:"cause_kind"`
	Participants []string      `json:"participants"`
}
type Ledger struct {
	Entries                    []Entry
	source, cause, participant map[string][]int
}

func Build(events []contracts.DomainEvent) Ledger {
	ordered := append([]contracts.DomainEvent(nil), events...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].Sequence < ordered[j].Sequence })
	l := Ledger{source: map[string][]int{}, cause: map[string][]int{}, participant: map[string][]int{}}
	for _, e := range ordered {
		var env operations.EventEnvelope
		_ = json.Unmarshal(e.PublicPayload, &env)
		p := env.Provenance
		source := p.SourceRef
		if source == "" {
			source = e.SourceRef
		}
		x := Entry{EventID: e.ID, Sequence: e.Sequence, Type: e.Type, SourceRef: source, OperationID: p.OperationID, CauseKind: p.CauseKind, Participants: append([]string(nil), p.Participants...)}
		i := len(l.Entries)
		l.Entries = append(l.Entries, x)
		if x.SourceRef != "" {
			l.source[x.SourceRef] = append(l.source[x.SourceRef], i)
		}
		if x.CauseKind != "" {
			l.cause[x.CauseKind] = append(l.cause[x.CauseKind], i)
		}
		for _, v := range x.Participants {
			l.participant[v] = append(l.participant[v], i)
		}
	}
	return l
}
func (l Ledger) BySource(v string) []Entry      { return l.pick(l.source[v]) }
func (l Ledger) ByCause(v string) []Entry       { return l.pick(l.cause[v]) }
func (l Ledger) ByParticipant(v string) []Entry { return l.pick(l.participant[v]) }
func (l Ledger) pick(xs []int) []Entry {
	r := make([]Entry, len(xs))
	for i, x := range xs {
		r[i] = l.Entries[x]
		r[i].Participants = append([]string(nil), r[i].Participants...)
	}
	return r
}
