package bgpfinder

const (
	ROUTEVIEWS = "routeviews"
)

// TODO: Finder implementation for the RouteViews archive
type RouteViewsFinder struct {
	// Cache of collectors
	collectors []Collector
}

func NewRouteViewsFinder() *RouteViewsFinder {
	f := &RouteViewsFinder{}

	return f
}

func (f *RouteViewsFinder) Projects() ([]string, error) {
	return []string{ROUTEVIEWS}, nil
}

func (f *RouteViewsFinder) Collectors(project string) ([]Collector, error) {
	if project != ROUTEVIEWS {
		return nil, nil
	}
	// TODO: turn this into a goroutine that periodically
	// refreshes collector list (and handles transient failures)?
	return f.getCollectors()
}

func (f *RouteViewsFinder) Find(query Query) ([]File, error) {
	// TODO
	return nil, nil
}

func (f *RouteViewsFinder) getCollectors() ([]Collector, error) {
	return nil, nil
}
