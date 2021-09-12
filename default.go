package bgpfinder

// Global finder instance that includes all the built-in finder
// implementations (RV and RIS for now).
//
// If you have a custom (private) finder, you can either register it
// with this finder instance, or use it directly.
var DefaultFinder = initDefaultFinder()

func initDefaultFinder() *MultiFinder {
	m := NewMultiFinder()
	// TODO: add each of our default finder providers (RV, RIS)
	return m
}

func Projects() []string {
	return DefaultFinder.Projects()
}

func Collectors(project string) []Collector {
	return DefaultFinder.Collectors(project)
}

func Find(query Query) []Result {
	return DefaultFinder.Find(query)
}
