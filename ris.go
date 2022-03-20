package bgpfinder

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/alistairking/bgpfinder/scraper"
)

const (
	RIS             = "ris"
	RIS_ARCHIVE_URL = "https://data.ris.ripe.net/"
	// XXX: it's tempting, but we can't use
	// https://www.ris.ripe.net/peerlist/ because it only lists
	// currently-active collectors.
	RIS_COLLECTORS_URL = RIS_ARCHIVE_URL
)

var (
	RIS_PROJECT = Project{Name: RIS}

	// RIS collectors are easy to find...
	risRRCPattern = regexp.MustCompile(`(rrc\d\d)`)
)

type RISFinder struct {
	// Cache of collectors
	mu            *sync.RWMutex
	collectors    []Collector
	collectorsErr error // set if collectors is nil, nil otherwise
}

func NewRISFinder() *RISFinder {
	f := &RISFinder{
		mu: &sync.RWMutex{},
	}

	// TODO: turn this into a goroutine that periodically
	// refreshes collector list (and handles transient failures)?
	c, err := f.getCollectors()
	f.collectors = c
	f.collectorsErr = err

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
	if project != "" && project != RIS {
		return nil, nil
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.collectors, f.collectorsErr
}

func (f *RISFinder) Collector(name string) (Collector, error) {
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
	return Collector{}, nil
}

func (f *RISFinder) Find(query Query) ([]File, error) {
	// TODO
	return nil, nil
}

func (f *RISFinder) getCollectors() ([]Collector, error) {
	links, err := scraper.ScrapeLinks(RIS_COLLECTORS_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to get collector list: %v", err)
	}
	// RRCs are easy to find.. perhaps a little too easy.. so let's throw
	// away dups
	seen := map[string]struct{}{}
	colls := []Collector{}
	for _, link := range links {
		m := risRRCPattern.FindStringSubmatch(link)
		if len(m) != 2 {
			continue
		}
		c := m[1]
		if _, exists := seen[c]; !exists {
			colls = append(colls, Collector{
				Project:      RIS_PROJECT,
				Name:         m[1],
				InternalName: m[1],
			})
			seen[c] = struct{}{}
		}
	}
	return colls, nil
}
