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

	// TODO: turn this into a goroutine that periodically
	// refreshes collector list (and handles transient failures)?
	c, _ := f.getCollectors()
	f.collectors = c

	return f
}

func (f *RISFinder) Projects() []string {
	return []string{RIS}
}

func (f *RISFinder) Collectors(project string) []Collector {
	if project != RIS {
		return nil
	}
	return f.collectors
}

func (f *RISFinder) Find(query Query) ([]Result, error) {
	// TODO
	return nil, nil
}

func (f *RISFinder) getCollectors() ([]Collector, error) {
	return nil, nil
}
