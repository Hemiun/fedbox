package user

type UserRequest struct {
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Tags     []string `json:"tags"`
	Comments string   `json:"comments"`
}
