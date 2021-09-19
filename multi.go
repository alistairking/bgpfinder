package bgpfinder

import (
	"fmt"
	"sync"
)

// Finder implementation that handles routing requests to a set of sub finder
// instances.
type MultiFinder struct {
	finders map[string]Finder
	mu      *sync.RWMutex
}

func NewMultiFinder(finders ...Finder) (*MultiFinder, error) {
	m := &MultiFinder{
		finders: map[string]Finder{},
		mu:      &sync.RWMutex{},
	}
	for _, f := range finders {
		err := m.AddFinder(f)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *MultiFinder) AddFinder(f Finder) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// TODO: handle project collisions
	projs, err := f.Projects()
	if err != nil {
		return err
	}
	for _, proj := range projs {
		m.finders[proj] = f
	}
	return nil

}

func (m *MultiFinder) Projects() ([]string, error) {
	// TODO: move this to AddFinder and cache result
	m.mu.RLock()
	defer m.mu.RUnlock()
	projs := make([]string, len(m.finders))
	idx := 0
	for p := range m.finders {
		projs[idx] = p
		idx++
	}
	return projs, nil
}

func (m *MultiFinder) Collectors(project string) ([]Collector, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if project != "" {
		proj, exists := m.finders[project]
		if !exists {
			// TODO: define our error types
			return nil, fmt.Errorf("Unknown project: '%s'", project)
		}
		return proj.Collectors(project)
	}
	allColls := []Collector{}
	for _, f := range m.finders {
		colls, err := f.Collectors(project)
		if err != nil {
			return nil, err
		}
		allColls = append(allColls, colls...)

	}
	return allColls, nil
}

func (m *MultiFinder) Find(query Query) ([]File, error) {
	// TODO
	return nil, nil
}
