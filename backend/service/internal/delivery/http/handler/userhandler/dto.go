package userhandler

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Surname  string `json:"surname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CurrentUserDTO struct {
	UserID  int64  `json:"userId"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Role    string `json:"role"`
}

type LoginResponse struct {
	Token       string         `json:"token"`
	CurrentUser CurrentUserDTO `json:"currentUser"`
}
