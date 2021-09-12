package bgpfinder

import (
	"time"
)

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
