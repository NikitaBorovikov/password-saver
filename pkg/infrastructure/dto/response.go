package dto

import "time"

type PasswordResponse struct {
	PasswordID int64  `json:"password_id"`
	Service    string `json:"service"`
	Password   string `json:"password"`
	Login      string `json:"login"`
}

type GetUserInfoResponse struct {
	UserID int64  `json:"userID"`
	Email  string `json:"email"`
}

type OKResponse struct {
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthCheckResponse struct {
	Status  string    `json:"status"`
	Details string    `json:"details"`
	Time    time.Time `json:"time"`
}

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
