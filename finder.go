package bgpfinder

import (
	"fmt"
	"strings"
	"time"
)

// Just a sketch of what the base Finder interface might look like.  Everything
// gets built on top of (or under, I guess) this.
type Finder interface {
	// Get the list of projects supported by this finder
	Projects() ([]Project, error)

	// Get a specific project
	Project(name string) (Project, error)

	// Get the list of collectors supported by the given project. All
	// projects if unset.
	Collectors(project string) ([]Collector, error)

	// Get a specific collector by name, returns the zero collector if there
	// is no such collector
	Collector(name string) (Collector, error)

	// Find all the BGP data URLs that match the given query
	Find(query Query) ([]File, error)
}

type Project struct {
	Name string `json:"name"`
}

func (p Project) String() string {
	return p.Name
}

func (p Project) AsCSV() string {
	return p.Name
}

type Collector struct {
	// Project name the collector belongs to
	Project Project `json:"project"`

	// Name of the collector
	Name string `json:"name"`

	// Project-internal name for this collector
	InternalName string `json:"internal_name"`
}

var ZeroCollector = Collector{}

func (c Collector) String() string {
	return fmt.Sprintf("%s:%s", c.Project, c.Name)
}

func (c Collector) AsCSV() string {
	return strings.Join([]string{
		c.Project.AsCSV(),
		c.Name,
		c.InternalName,
	}, ",")
}

// TODO: add BGPStream backwards compat names.

//go:generate enumer -type=DumpType -json -text -linecomment
type DumpType uint8

const (
	DUMP_TYPE_ANY     DumpType = 0 // any
	DUMP_TYPE_RIB     DumpType = 1 // rib
	DUMP_TYPE_UPDATES DumpType = 2 // updates
)

// TODO: think about how this should work -- just keep it simple! no
// complex query structures
//
// TODO: add Validate method (e.g., From is before Until, IsADumpType, etc.)
type Query struct {
	// Collectors to search for. All collectors if unset/empty
	Collectors []Collector

	// Query window start time (inclusive)
	From time.Time

	// Query window end time (exclusive)
	Until time.Time

	// Dump type to search for. Any type if unset
	DumpType DumpType
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

	// Dump type
	DumpType DumpType

	// TODO: other things? (file size?)
}
