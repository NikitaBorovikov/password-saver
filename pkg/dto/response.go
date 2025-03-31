package dto

type (
	PasswordResponse struct {
		PasswordID int64  `json:"password_id"`
		Service    string `json:"service"`
		Password   string `json:"password"`
		Login      string `json:"login"`
	}

	GetUserInfoResponse struct {
		UserID int64  `json:"userID"`
		Email  string `json:"email"`
	}

	OKResponse struct {
		Data interface{} `json:"data"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

func NewOKResponse(data interface{}) OKResponse {
	return OKResponse{
		Data: data,
	}
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: err.Error(),
	}
}

func NewGetUserInfoResponse(userID int64, email string) *GetUserInfoResponse {
	return &GetUserInfoResponse{
		UserID: userID,
		Email:  email,
	}
}
