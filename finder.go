package bgpfinder

import (
	"net/url"
)

// Just a sketch of what the base Finder interface might look like.
// Everything gets built on top of (or under, I guess) this.
type Finder interface {
	// Get the list of projects supported by this finder
	Projects() []string

	// Get the list of collectors supported by the given
	// project. All projects if unset.
	Collectors(project string) []Collector

	// Find all the BGP data URLs that match the given query
	Find(query Query) []Result
}

// Represents a single BGP file found by a Finder.
type Result struct {
	// URL of the file
	URL url.URL

	// Collector that collected this file
	Collector Collector

	// TODO: other things?
}
