package inducedates

// Query query with ott_ids
type Query struct {
	Ottids []int `json:"ott_ids"`
}

// Newick return newick string json
type Newick struct {
	NewString string `json:"newick"`
	Unmatched []int  `json:"unmatched_ott_ids"`
}

// NewickQuery in newick query
type NewickQuery struct {
	NewIn string `json:"newick"`
}
