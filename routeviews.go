package bgpfinder

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alistairking/bgpfinder/scraper"
)

const (
	ROUTEVIEWS                = "routeviews"
	ROUTEVIEWS_ARCHIVE_URL    = "http://archive.routeviews.org/"
	ROUTEVIEWS_COLLECTORS_URL = ROUTEVIEWS_ARCHIVE_URL
)

var (
	// These are last-resort overrides to "fix" an out-of-pattern RV
	// collector name.
	ROUTEVIEWS_COLLECTOR_OVERRIDES = map[string][2]string{
		"": {"route-views2", "rv2"},
	}

	// Various collector-name parsing regexes
	// XXX: could collapse these into a single regex perhaps
	rvCollDigitsOnly = regexp.MustCompile(`^route-views(\d+)$`)
	rvCollName       = regexp.MustCompile(`^route-views(\d+)?\.([a-zA-Z0-9]+)$`)
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
		return nil, fmt.Errorf("failed to get collector list: %v", err)
	}
	// we're expecting to get things like:
	// '/route-views.nwax/bgpdata'
	// but there are odd things like:
	// '/bgpdata'
	// and the telnet links like:
	// 'telnet://route-views.perth.routeviews.org'
	colls := []Collector{}
	for _, link := range links {
		origLink := link
		intName := ""
		if !strings.HasSuffix(link, "/bgpdata") {
			continue
		}
		link = strings.TrimSuffix(link, "/bgpdata")
		link = strings.TrimPrefix(link, "/")
		// now we're left with three classes of collector (and one oddball)
		// '<DIGIT>' => 'rv<DIGIT>'
		m := rvCollDigitsOnly.FindStringSubmatch(link)
		if len(m) == 2 {
			intName = link
			link = "rv" + m[1]
		}
		m = rvCollName.FindStringSubmatch(link)
		if m != nil {
			intName = link
			// route-views.sg
			link = m[2]
			if len(m) == 3 {
				// route-views2.saopaulo
				link += m[1]
			}
		}
		// 'route-views.<NAME> => '<NAME>'
		// 'route-views<DIGIT>.<NAME> => '<NAME><DIGIT>'
		override, exists := ROUTEVIEWS_COLLECTOR_OVERRIDES[link]
		if exists {
			intName = override[0]
			link = override[1]
		}
		if intName == "" {
			return nil, fmt.Errorf("unexpected collector pattern: '%s' ('%s'). "+
				"Please file a parser bug report at "+
				"https://github.com/alistairking/bgpfinder/issues", link, origLink)
		}
		colls = append(colls, Collector{
			Project:      ROUTEVIEWS,
			Name:         link,
			InternalName: intName,
		})
	}
	return colls, nil
}
