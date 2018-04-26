package inducedates

// Query query with ott_ids
type Query struct {
	Ottids []int `json:"ott_ids"`
}

// Newick return newick string json
type Newick struct {
	NewString string `json:"newick"`
}