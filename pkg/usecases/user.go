package usecases

import (
	"errors"
	"password-saver/pkg/dto"
	apperrors "password-saver/pkg/errors"
	"password-saver/pkg/logs"
	"password-saver/pkg/model"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepository model.UserRepository
}

func NewUserUseCase(ur model.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepository: ur,
	}
}

func (uc *UserUseCase) Registration(req *dto.AuthRequest) (int64, error) {

	if err := validateAuthRequest(req); err != nil {
		logrus.Errorf(logs.FailedToValidateUser, err)
		return 0, err
	}

	hashPassword, err := hashPassword(req.Password)
	if err != nil {
		logrus.Errorf(logs.FailedToHashPassword, err)
		return 0, apperrors.ErrHashPassword
	}

	regDate := getTodayDate()

	user := newUser(0, req.Email, hashPassword, regDate)

	userID, err := uc.UserRepository.Registration(user)
	if err != nil {
		return 0, handleUserRepositoryError(err, req.Email)
	}

	return userID, nil
}

func (uc *UserUseCase) LogIn(req *dto.AuthRequest) (*model.User, error) {

	if err := validateAuthRequest(req); err != nil {
		logrus.Errorf(logs.FailedToValidateUser, err)
		return nil, err
	}

	user, err := uc.UserRepository.LogIn(req)
	if err != nil {
		return nil, handleUserRepositoryError(err, req.Email)
	}

	if !comparePassword(req.Password, user.HashPassword) {
		logrus.Errorf(logs.FailedToComparePasswords, err)
		return nil, apperrors.ErrComparePasswords
	}

	sanitizeUserStruct(user)

	return user, nil
}

func (uc *UserUseCase) Update(req *dto.UpdateUserRequest, userID int64) error {

	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		return handleUserRepositoryError(err, userID)
	}

	if err := validateUpdateRequest(req); err != nil {
		logrus.Errorf(logs.FailedToValidateUser, err)
		return err
	}

	if !comparePassword(req.OldPassword, user.HashPassword) {
		logrus.Errorf(logs.FailedToComparePasswords, err)
		return apperrors.ErrComparePasswords
	}

	user.HashPassword, err = hashPassword(req.NewPassword)
	if err != nil {
		logrus.Errorf(logs.FailedToHashPassword, err)
		return apperrors.ErrHashPassword
	}

	if err := uc.UserRepository.Update(user); err != nil {
		return handleUserRepositoryError(err, userID)
	}

	return nil
}

func (uc *UserUseCase) GetByID(userID int64) (*model.User, error) {
	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		return nil, handleUserRepositoryError(err, userID)
	}
	sanitizeUserStruct(user)

	return user, nil
}

func (uc *UserUseCase) Delete(userID int64) error {
	if err := uc.UserRepository.Delete(userID); err != nil {
		return handleUserRepositoryError(err, userID)
	}

	return nil
}

func validateAuthRequest(req *dto.AuthRequest) error {
	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		return handleValidateAuthErrors(err)
	}
	return nil
}

func validateUpdateRequest(req *dto.UpdateUserRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return handleValidateUpdateUserError(err)
	}

	return nil
}

func handleValidateAuthErrors(err error) error {
	var validateErrs validator.ValidationErrors

	if !errors.As(err, &validateErrs) {
		return apperrors.ErrValidateUser
	}

	for _, e := range validateErrs {

		switch e.Field() {
		case "Email":
			return apperrors.ErrValidateEmailField
		case "Password":
			return apperrors.ErrValidateUserPasswordField
		}
	}
	return apperrors.ErrValidateUser
}

func handleValidateUpdateUserError(err error) error {
	var validateErrs validator.ValidationErrors

	if !errors.As(err, &validateErrs) {
		return apperrors.ErrValidateUser
	}

	for _, e := range validateErrs {

		switch e.Field() {
		case "OldPassword":
			return apperrors.ErrValidateOldPasswordField
		case "NewPassword":
			return apperrors.ErrValidateNewPasswordField
		}
	}
	return apperrors.ErrValidateUser
}

func getTodayDate() string {
	return time.Now().Format(time.RFC3339)
}

func hashPassword(inputPassword string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}

func comparePassword(inputPassword, hashPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(inputPassword)); err != nil {
		return false
	}
	return true
}

func sanitizeUserStruct(u *model.User) {
	u.HashPassword = ""
}

func newUser(userID int64, email, hashPassword, regDate string) *model.User {
	return &model.User{
		UserID:       userID,
		Email:        email,
		HashPassword: hashPassword,
		RegDate:      regDate,
	}
}

func handleUserRepositoryError(err error, data interface{}) error {
	switch err {
	case apperrors.ErrUserNotFound:
		logrus.Errorf("userID: %d %v", data, err)
		return apperrors.ErrUserNotFound

	case apperrors.ErrDuplicateUser:
		logrus.Errorf("email: %s %v", data, err)
		return apperrors.ErrDuplicateUser

	default:
		logrus.Errorf(logs.InternalDBError, err)
		return apperrors.ErrDatabaseInternal
	}
}
