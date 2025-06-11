package requests

type RegisterRequest struct {
	Fullname string  `json:"fullname" validate:"required,min=3"`
	Username string  `json:"username" validate:"required,min=3"`
	Password string  `json:"password" validate:"required,min=6"`
	Role     string  `json:"role" validate:"required,oneof=admin employee"`
	Salary   float64 `json:"salary" validate:"required"`
}
