package web

type UserUpdateRequest struct {
	Id       string `validate:"required,uuid4" json:"id"`
	Name     string `validate:"required,min=3" json:"name"`
	Email    string `validate:"required,email" json:"email"`
	Password string `json:"password"`
}
