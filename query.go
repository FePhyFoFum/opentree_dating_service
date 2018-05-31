package inducedates

// Query query with ott_ids
type Query struct {
	Ottids []int `json:"ott_ids"`
}

// GBIFQuery for getting ottids from gbifids
type GBIFQuery struct {
	GBIFids []int `json:"gbif_ids"`
}

// OttidResults for returning ottid results
type OttidResults struct {
	Ottids    []int `json:"ott_ids"`
	Unmatched []int `json:"unmatched_gbif_ids"`
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
