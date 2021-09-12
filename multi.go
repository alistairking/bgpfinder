package bgpfinder

import (
	"sync"
)

// Finder implementation that handles routing requests to a set of sub
// finder instances.
type MultiFinder struct {
	finders map[string]Finder
	mu      *sync.RWMutex
}

func NewMultiFinder(finders ...Finder) *MultiFinder {
	m := &MultiFinder{
		finders: map[string]Finder{},
		mu:      &sync.RWMutex{},
	}
	for _, f := range finders {
		m.AddFinder(f)
	}
	return m
}

func (m *MultiFinder) AddFinder(f Finder) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO: handle project collisions
	for _, proj := range f.Projects() {
		m.finders[proj] = f
	}

}

func (m *MultiFinder) Projects() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	projs := make([]string, len(m.finders))
	idx := 0
	for p := range m.finders {
		projs[idx] = p
		idx++
	}
	return projs
}

func (m *MultiFinder) Collectors(project string) []Collector {
	// TODO
	return nil
}

func (m *MultiFinder) Find(query Query) ([]Result, error) {
	// TODO
	return nil, nil
}
