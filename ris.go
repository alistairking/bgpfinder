package bgpfinder

const (
	RIS = "ris"
)

var (
	RIS_PROJECT = Project{Name: RIS}
)

// TODO: add Finder implementation for the RIS archive
type RISFinder struct {
}

func NewRISFinder() *RISFinder {
	f := &RISFinder{}
	return f
}

func (f *RISFinder) Projects() ([]Project, error) {
	return []Project{RIS_PROJECT}, nil
}

func (f *RISFinder) Project(name string) (Project, error) {
	if name == "" || name == RIS {
		return RIS_PROJECT, nil
	}
	return Project{}, nil
}

func (f *RISFinder) Collectors(project string) ([]Collector, error) {
	if project != RIS {
		return nil, nil
	}
	// TODO
	return nil, nil
}

func (f *RISFinder) Collector(name string) (Collector, error) {
	// TODO
	return Collector{}, nil
}

func (f *RISFinder) Find(query Query) ([]File, error) {
	// TODO
	return nil, nil
}
