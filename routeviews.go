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
	if f.collectors != nil {
		return f.collectors, nil
	}
	c, err := f.getCollectors()
	if err != nil {
		return nil, err
	}
	f.collectors = c
	return c, nil
}

func (f *RouteViewsFinder) Find(query Query) ([]File, error) {
	// TODO
	return nil, nil
}

func (f *RouteViewsFinder) getCollectors() ([]Collector, error) {
	// If we can find a Go rsync client (not a wrapper) we could
	// just do `rsync archive.routeviews.org::` and do some light
	// parsing on the output.
	//
	// Alternatively we can parse http://archive.routeviews.org/
	// and look for links like
	// http://archive.routeviews.org/route-views.chicago/bgpdata
	//
	// I'd like to not repeat the original mistake we made and
	// would prefer to call the above collector "chicago" rather
	// than "route-views.chicago". We can always map these back if
	// we need to for BGPStream compat.
	return nil, nil
}
