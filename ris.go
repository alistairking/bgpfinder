package bgpfinder

const (
	RIS = "ris"
)

// TODO: add Finder implementation for the RIS archive
type RISFinder struct {
	// Cache of collectors
	collectors []Collector
}

func NewRISFinder() *RISFinder {
	f := &RISFinder{}
	return f
}

func (f *RISFinder) Projects() ([]string, error) {
	return []string{RIS}, nil
}

func (f *RISFinder) Collectors(project string) ([]Collector, error) {
	if project != RIS {
		return nil, nil
	}
	// TODO: turn this into a goroutine that periodically
	// refreshes collector list (and handles transient failures)?
	return f.getCollectors()
}

func (f *RISFinder) Find(query Query) ([]File, error) {
	// TODO
	return nil, nil
}

func (f *RISFinder) getCollectors() ([]Collector, error) {
	return nil, nil
}
