package bgpfinder

import (
	"fmt"
	"strings"

	"github.com/alistairking/bgpfinder/scraper"
)

const (
	ROUTEVIEWS                = "routeviews"
	ROUTEVIEWS_ARCHIVE_URL    = "http://archive.routeviews.org/"
	ROUTEVIEWS_COLLECTORS_URL = ROUTEVIEWS_ARCHIVE_URL
)

var (
	ROUTEVIEWS_COLLECTOR_OVERRIDES = map[string]string{
		"":           "rv2",
		"3":          "rv3",
		"4":          "rv4",
		"6":          "rv6",
		"2.saopaulo": "saopaulo2",
	}
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
	if project != "" && project != ROUTEVIEWS {
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
	// If we can find a Go rsync client (not a wrapper) we could just do
	// `rsync archive.routeviews.org::` and do some light parsing on the
	// output.
	//
	// Alternatively we can parse http://archive.routeviews.org/ and look
	// for links like
	// http://archive.routeviews.org/route-views.chicago/bgpdata
	//
	// I'd like to not repeat the original mistake we made and would prefer
	// to call the above collector "chicago" rather than
	// "route-views.chicago". We can always map these back if we need to for
	// BGPStream compat.
	links, err := scraper.ScrapeLinks(ROUTEVIEWS_COLLECTORS_URL)
	if err != nil {
		return nil, fmt.Errorf("Failed to get collector list: %v", err)
	}
	// we're expecting to get things like:
	// '/route-views.nwax/bgpdata'
	// but there are odd things like:
	// '/bgpdata'
	// and the telnet links like:
	// 'telnet://route-views.perth.routeviews.org'
	colls := []Collector{}
	for _, link := range links {
		// TODO: make this more generic
		if !strings.HasSuffix(link, "/bgpdata") {
			continue
		}
		link = strings.TrimPrefix(link, "/route-views")
		link = strings.TrimSuffix(link, "/bgpdata")
		link = strings.TrimPrefix(link, ".")
		// TODO:...
		override, exists := ROUTEVIEWS_COLLECTOR_OVERRIDES[link]
		if exists {
			link = override
		}
		colls = append(colls, Collector{
			Project: ROUTEVIEWS,
			Name:    link,
		})
	}
	return colls, nil
}
