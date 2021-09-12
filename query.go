package bgpfinder

import (
	"time"
)

// TODO: think about how this should work -- just keep it simple! no
// complex query structures
type Query struct {
	// Name of the project. All projects if unset
	Project string

	// Name of the collector. All collectors if unset
	Collector string

	// Query window start time (inclusive)
	From time.Time

	// Query window end time (exclusive)
	Until time.Time
}
