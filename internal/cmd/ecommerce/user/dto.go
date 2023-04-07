package user

// UserDTO is a dto for view/create/update a user
type UserDTO struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Tags     []string `json:"tags"`
	Comments string   `json:"comments"`
}
