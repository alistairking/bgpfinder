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

	// TODO: turn this into a goroutine that periodically
	// refreshes collector list (and handles transient failures)?
	c, _ := f.getCollectors()
	f.collectors = c

	return f
}

func (f *RouteViewsFinder) Projects() []string {
	return []string{ROUTEVIEWS}
}

func (f *RouteViewsFinder) Collectors(project string) []Collector {
	if project != ROUTEVIEWS {
		return nil
	}
	return f.collectors
}

func (f *RouteViewsFinder) Find(query Query) ([]Result, error) {
	// TODO
	return nil, nil
}

func (f *RouteViewsFinder) getCollectors() ([]Collector, error) {
	return nil, nil
}
