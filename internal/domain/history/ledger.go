// Package history builds deterministic, rebuildable provenance indexes from events.
package history

import (
	"encoding/json"
	"github.com/NordCoder/unmatched-web/internal/domain/contracts"
	"github.com/NordCoder/unmatched-web/internal/domain/model"
	"sort"
)

type Entry struct {
	EventID                    model.EventID
	Sequence                   uint64
	Type, SourceRef, CauseKind string
	Participants               []string
}
type Ledger struct {
	Entries                    []Entry
	source, cause, participant map[string][]int
}

func Build(events []contracts.DomainEvent) Ledger {
	sort.Slice(events, func(i, j int) bool { return events[i].Sequence < events[j].Sequence })
	l := Ledger{source: map[string][]int{}, cause: map[string][]int{}, participant: map[string][]int{}}
	for _, e := range events {
		var p struct {
			CauseKind    string   `json:"cause_kind"`
			Participants []string `json:"participants"`
		}
		_ = json.Unmarshal(e.PublicPayload, &p)
		x := Entry{e.ID, e.Sequence, e.Type, e.SourceRef, p.CauseKind, p.Participants}
		i := len(l.Entries)
		l.Entries = append(l.Entries, x)
		l.source[x.SourceRef] = append(l.source[x.SourceRef], i)
		l.cause[x.CauseKind] = append(l.cause[x.CauseKind], i)
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
	}
	return r
}
