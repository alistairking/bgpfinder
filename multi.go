package bgpfinder

// Finder implementation that handles routing requests to a set of sub
// finder instances.
type MultiFinder struct {
	// TODO
}

func NewMultiFinder() *MultiFinder {
	// TODO
	m := &MultiFinder{}
	return m
}

func (m *MultiFinder) AddFinder(f Finder) error {
	// TODO
	return nil
}

func (m *MultiFinder) Projects() []string {
	// TODO
	return nil
}

func (m *MultiFinder) Collectors(project string) []Collector {
	// TODO
	return nil
}

func (m *MultiFinder) Find(query Query) []Result {
	// TODO
	return nil
}
