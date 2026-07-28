package runtime

import (
	"fmt"
	"sync"
)

// IDProvider allocates opaque runtime identity. Tests can inject a deterministic
// implementation; production can later supply UUID-backed allocation.
type IDProvider interface {
	Next(kind string) string
}

type SequenceIDProvider struct {
	mu      sync.Mutex
	prefix  string
	counter map[string]uint64
}

func NewSequenceIDProvider(prefix string) *SequenceIDProvider {
	return &SequenceIDProvider{prefix: prefix, counter: make(map[string]uint64)}
}

func (p *SequenceIDProvider) Next(kind string) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.counter[kind]++
	if p.prefix == "" {
		return fmt.Sprintf("%s-%d", kind, p.counter[kind])
	}
	return fmt.Sprintf("%s-%s-%d", p.prefix, kind, p.counter[kind])
}
