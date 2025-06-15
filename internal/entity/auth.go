package entity

type User struct {
	Name     string
	Password string
	Token    string
}

type LoginResponse struct {
	Token string `json:"token" validate:"required"`
}
