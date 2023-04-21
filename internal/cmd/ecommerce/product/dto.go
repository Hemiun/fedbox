package product

// Product is a dto for view/create/update a product
type Product struct {
	Id      string         `json:"id"`
	Name    string         `json:"name"`
	Summary string         `json:"summary"`
	Content ProductContent `json:"content"`
	Tags    []string       `json:"tags"`
}

// ProductContent contains custom (not ActivityPub) properties
type ProductContent struct {
	Price string `json:"price"`
	Color string `json:"color"`
	Size  string `json:"size"`
}
