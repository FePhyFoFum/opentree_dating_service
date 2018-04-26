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
		"InduceSubtree",
		"POST",
		"/induce_subtree",
		InduceSubtree,
	},
}
