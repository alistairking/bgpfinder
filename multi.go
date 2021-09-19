package bgpfinder

import (
	"fmt"
	"sync"
)

// Finder implementation that handles routing requests to a set of sub finder
// instances.
type MultiFinder struct {
	finders     map[string]Finder
	allProjects []Project
	projects    map[string]Project
	projColls   map[string][]Collector
	mu          *sync.RWMutex
}

func NewMultiFinder(finders ...Finder) (*MultiFinder, error) {
	m := &MultiFinder{
		finders:     map[string]Finder{},
		allProjects: []Project{},
		projects:    map[string]Project{},
		projColls:   map[string][]Collector{}, // lazy-loaded
		mu:          &sync.RWMutex{},
	}
	for _, f := range finders {
		err := m.AddFinder(f)
		if err != nil {
			return nil, err
		}
	}
	// TODO: kick off a slow collector updater goroutine?
	return m, nil
}

func (m *MultiFinder) AddFinder(f Finder) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	projs, err := f.Projects()
	if err != nil {
		return err
	}
	for _, proj := range projs {
		name := proj.Name
		_, collision := m.projects[name]
		if collision {
			return fmt.Errorf("project already added: %s", name)
		}
		m.finders[name] = f
		m.allProjects = append(m.allProjects, proj)
		m.projects[name] = proj
		// leave projColls to be loaded when it's needed
	}
	return nil
}

func (m *MultiFinder) Projects() ([]Project, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.allProjects, nil
}

func (m *MultiFinder) Project(name string) (Project, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	proj, exists := m.projects[name]
	if !exists {
		// TODO: define our error types
		return Project{}, fmt.Errorf("unknown project: '%s'", name)
	}
	return proj, nil
}

func (m *MultiFinder) Collectors(project string) ([]Collector, error) {
	if project != "" {
		f, exists := m.getFinderByProject(project)
		if !exists {
			return nil, fmt.Errorf("unknown project: '%s'", project)
		}
		colls, err := f.Collectors(project)
		if err != nil {
			return nil, err
		}
		m.mu.Lock()
		defer m.mu.Unlock()
		m.projColls[project] = colls
		return colls, nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	allColls := []Collector{}
	for _, f := range m.finders {
		colls, err := f.Collectors(project)
		if err != nil {
			return nil, err
		}
		m.projColls[project] = colls //< needs write lock
		allColls = append(allColls, colls...)
	}
	return allColls, nil
}

func (m *MultiFinder) Collector(name string) (Collector, error) {
	// tricky, we don't know where to send this request.
	// TODO: we should cache project->collector mappings
	colls, err := m.Collectors("")
	if err != nil {
		return Collector{}, err
	}
	for _, coll := range colls {
		if coll.Name == name {
			return coll, nil
		}
	}
	return Collector{}, nil
}

func (m *MultiFinder) Find(query Query) ([]File, error) {
	// TODO
	return nil, nil
}

func (m *MultiFinder) getFinderByProject(projName string) (Finder, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, exists := m.finders[projName]
	return f, exists
}
