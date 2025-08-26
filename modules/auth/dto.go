package auth

type RegisteUserDto struct {
	Username        string `json:"username" binding:"required,alphanum" validate:"required,alphanum"`
	Email           string `json:"email" binding:"required,email" validate:"required,email"`
	Password        string `json:"password" binding:"required,min=8" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" validate:"required,eqfield=Password"`
}

type LoginUserDto struct {
	UsernameEmail string `json:"username_email" binding:"required" validate:"required"`
	Password      string `json:"password" binding:"required" validate:"required"`
}
