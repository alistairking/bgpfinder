package bgpfinder

import (
	"fmt"
	"time"
)

// Just a sketch of what the base Finder interface might look like.
// Everything gets built on top of (or under, I guess) this.
type Finder interface {
	// Get the list of projects supported by this finder
	Projects() ([]string, error)

	// Get the list of collectors supported by the given
	// project. All projects if unset.
	Collectors(project string) ([]Collector, error)

	// Find all the BGP data URLs that match the given query
	Find(query Query) ([]File, error)
}

type Collector struct {
	// Project name the collector belongs to
	Project string `json:"project"`

	// Name of the collector
	Name string `json:"name"`
}

func (c Collector) String() string {
	return fmt.Sprintf("%s:%s", c.Project, c.Name)
}

// TODO: add BGPStream backwards compat names. This may neccesitate
// turning Project into a special type

// TODO: think about how this should work -- just keep it simple! no
// complex query structures
type Query struct {
	// Collectors to search for. All collectors if unset/empty
	Collector []Collector

	// Query window start time (inclusive)
	From time.Time

	// Query window end time (exclusive)
	Until time.Time
}

// Represents a single BGP file found by a Finder.
// TODO: better name for this. Dump is a candidate.
type File struct {
	// URL of the file
	URL string

	// Collector that collected this file
	Collector Collector

	// Nominal dump duration
	Duration time.Duration

	// TODO: other things? (file size?)
}
