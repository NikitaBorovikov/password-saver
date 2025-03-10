package dto

type (
	RegRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"min=7,max=40"`
	}

	LogInRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	UpdateUserRequest struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"min=7,max=40"`
	}

	PasswordRequest struct {
		Service  string `json:"service" validate:"min=1,max=100"`
		Password string `json:"password" validate:"min=1,max=100"`
	}
)
