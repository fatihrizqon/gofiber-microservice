package request

type LoginRequest struct {
	Email    string `validate:"required,email,min=1,max=254" json:"email" example:"johndoe@example.com"`
	Password string `validate:"required,min=8" json:"password" example:"yoursecretpassword"`
}

type RegisterRequest struct {
	Username string `validate:"required,min=1,max=20" json:"username"`
	Email    string `validate:"required,email,min=1,max=254" json:"email"`
	Password string `validate:"required,min=8" json:"password"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}
