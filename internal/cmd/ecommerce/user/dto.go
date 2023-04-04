package user

// UserRequest is dto for create user request
type UserRequest struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Tags     []string `json:"tags"`
	Comments string   `json:"comments"`
}
