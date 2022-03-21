package bgpfinder

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/alistairking/bgpfinder/scraper"
)

const (
	ROUTEVIEWS                = "routeviews"
	ROUTEVIEWS_ARCHIVE_URL    = "http://archive.routeviews.org/"
	ROUTEVIEWS_COLLECTORS_URL = ROUTEVIEWS_ARCHIVE_URL
)

var (
	ROUTEVIEWS_PROJECT = Project{Name: ROUTEVIEWS}

	// These are last-resort overrides to "fix" an out-of-pattern RV
	// collector name.
	ROUTEVIEWS_COLLECTOR_OVERRIDES = map[string][2]string{
		"": {"route-views2", "rv2"},
	}

	ROUTEVIEWS_DUMP_TYPES = map[DumpType]string{
		DUMP_TYPE_RIB:     "RIBS",
		DUMP_TYPE_UPDATES: "UPDATES",
	}

	// Various collector-name parsing regexes
	// XXX: could collapse these into a single regex perhaps
	rvCollDigitsOnly = regexp.MustCompile(`^route-views(\d+)$`)
	rvCollName       = regexp.MustCompile(`^route-views(\d+)?\.([a-zA-Z0-9]+)$`)
)

// TODO: Finder implementation for the RouteViews archive
// TODO: refactor a this common caching-finder code out so that RIS and PCH can use it
type RouteViewsFinder struct {
	// Cache of collectors
	mu            *sync.RWMutex
	collectors    []Collector
	collectorsErr error // set if collectors is nil, nil otherwise
}

func NewRouteViewsFinder() *RouteViewsFinder {
	f := &RouteViewsFinder{
		mu: &sync.RWMutex{},
	}

	// TODO: turn this into a goroutine that periodically
	// refreshes collector list (and handles transient failures)?
	c, err := f.getCollectors()
	f.collectors = c
	f.collectorsErr = err

	return f
}

func (f *RouteViewsFinder) Projects() ([]Project, error) {
	return []Project{ROUTEVIEWS_PROJECT}, nil
}

func (f *RouteViewsFinder) Project(name string) (Project, error) {
	if name == "" || name == ROUTEVIEWS {
		return ROUTEVIEWS_PROJECT, nil
	}
	return Project{}, nil
}

func (f *RouteViewsFinder) Collectors(project string) ([]Collector, error) {
	if project != "" && project != ROUTEVIEWS {
		return nil, nil
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.collectors, f.collectorsErr
}

func (f *RouteViewsFinder) Collector(name string) (Collector, error) {
	if f.collectorsErr != nil {
		return Collector{}, f.collectorsErr
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	// TODO: add a map to avoid the linear search
	for _, c := range f.collectors {
		if c.Name == name {
			return c, nil
		}
	}
	// not found
	return Collector{}, nil
}

func (f *RouteViewsFinder) Find(query Query) ([]File, error) {
	// ok, let's figure out which collectors we should query
	if len(query.Collectors) == 0 {
		// give them everything we got
		c, err := f.Collectors("")
		if err != nil {
			return nil, err
		}
		query.Collectors = c
	}
	results := []File{}
	// RV archives data by collector, so we want to do a collector-first
	// search
	for _, coll := range query.Collectors {
		// we can't trust anything other than the collector
		// name, so let's fix that
		fColl, err := f.Collector(coll.Name)
		if err != nil {
			return nil, err
		}
		if fColl == ZeroCollector {
			// invalid collector in query
			return nil, fmt.Errorf("Invalid collector: %+v", coll)
		}
		cRes, err := f.findFiles(fColl, query)
		if err != nil {
			// TODO: probably don't need to give up the whole
			// search...
			return nil, err
		}
		if cRes != nil {
			results = append(results, cRes...)
		}
	}
	return results, nil
}

func (f *RouteViewsFinder) getCollectors() ([]Collector, error) {
	// If we could find a Go rsync client (not a wrapper) we could just do
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
			Project:      ROUTEVIEWS_PROJECT,
			Name:         link,
			InternalName: intName,
		})
	}
	return colls, nil
}

func (f *RouteViewsFinder) collURL(coll Collector) string {
	// RV2 is the only special-case URL (afaik)
	if coll.Name == "rv2" {
		return ROUTEVIEWS_ARCHIVE_URL + "bgpdata/"
	}
	return ROUTEVIEWS_ARCHIVE_URL + coll.InternalName + "/bgpdata/"
}

func (f *RouteViewsFinder) monthURL(coll Collector, month time.Time) string {
	return f.collURL(coll) + month.Format("2006.01") + "/"
}

func (f *RouteViewsFinder) dumpTypeURL(coll Collector, month time.Time, rvdt string) string {
	return f.monthURL(coll, month) + rvdt + "/"
}

func (f *RouteViewsFinder) dumpTypes(dt DumpType) ([]string, error) {
	if dt == DUMP_TYPE_ANY {
		all := []string{}
		for _, rvdt := range ROUTEVIEWS_DUMP_TYPES {
			all = append(all, rvdt)
		}
		return all, nil
	}
	rvt, ok := ROUTEVIEWS_DUMP_TYPES[dt]
	if !ok {
		return nil, fmt.Errorf("invalid RouteViews dump type: %v", dt)
	}
	return []string{rvt}, nil
}

func (f *RouteViewsFinder) findFiles(coll Collector, query Query) ([]File, error) {
	// RV archive is organized by YYYY.MM, so we first iterate
	// over the months in our query range (there has to be at
	// least one)
	//
	// But first, let's figure out our dump type(s)
	rvdts, err := f.dumpTypes(query.DumpType)
	if err != nil {
		return nil, err
	}

	res := []File{}
	cur := query.From
	for cur.Before(query.Until) {
		for _, rvdt := range rvdts {
			dtUrl := f.dumpTypeURL(coll, cur, rvdt)
			if dtRes, err := f.findFilesForURL(res, dtUrl, query); err != nil {
				return nil, err
			} else {
				res = dtRes
			}

		}
		cur = cur.AddDate(0, 1, 0)
	}
	return nil, nil
}

func (f *RouteViewsFinder) findFilesForURL(res []File, url string, query Query) ([]File, error) {
	// Here we have something like:
	// url=http://archive.routeviews.org/route-views3/bgpdata/2020.09/RIBS/
	// now we need to grab the files there and figure out which
	// ones actually match our query.

	links, err := scraper.ScrapeLinks(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get file list from %s: %v", url, err)
	}

	fmt.Println(url, len(links), links[50])

	return res, nil
}
