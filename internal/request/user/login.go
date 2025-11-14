package user

type Login struct {
	Phone string `json:"phone" binding:"required,phone"`
	Code  string `json:"code" binding:"required"`
}
