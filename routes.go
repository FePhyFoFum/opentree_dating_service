package inducedates

import "net/http"

// Route basic route type
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes list of routes
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Emot",
		"GET",
		"/emot",
		Emot,
	},
	Route{
		"InducedSubtree",
		"POST",
		"/induced_subtree",
		InducedSubtree,
	},
	Route{
		"RenameTree",
		"POST",
		"/rename_tree",
		RenameTree,
	},
	Route{
		"RenameTreeNCBI",
		"POST",
		"/rename_tree_ncbi",
		RenameTreeNCBI,
	},
}
