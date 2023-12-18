package web

type UserCreateRequest struct {
	Id        string `validate:"required,uuid4" json:"id"`
	Name      string `validate:"required,min=3" json:"name"`
	Email     string `validate:"required,email" json:"email"`
	Password  string `validate:"required" json:"password"`
	CreatedAt int64  `validate:"required" json:"created_at"`
}
