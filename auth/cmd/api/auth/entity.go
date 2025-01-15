package auth

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}
