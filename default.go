package bgpfinder

// Global finder instance that includes all the built-in finder
// implementations (RV and RIS for now).
//
// If you have a custom (private) finder, you can either register it
// with this finder instance, or use it directly.
var DefaultFinder = NewMultiFinder(
	NewRouteViewsFinder(),
	NewRISFinder(),
)

func Projects() []string {
	return DefaultFinder.Projects()
}

func Collectors(project string) []Collector {
	return DefaultFinder.Collectors(project)
}

func Find(query Query) ([]Result, error) {
	return DefaultFinder.Find(query)
}
