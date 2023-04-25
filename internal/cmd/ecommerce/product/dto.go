package product

// Product is a dto for view/create/update a product
type Product struct {
	Id      string         `json:"id"`
	Name    string         `json:"name"`
	Summary string         `json:"summary"`
	Content map[string]any `json:"content"`
	Tags    []string       `json:"tags"`
	Images  []Image        `json:"images"`
}

type Image struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	URL     string `json:"url"`
}
