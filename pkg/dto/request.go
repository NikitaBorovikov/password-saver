package dto

type (
	RegRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LogInRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	UpdateUserRequest struct {
		UserID      int64  `json:"-"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
)
