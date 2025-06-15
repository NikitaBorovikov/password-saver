package usecases

import (
	"password-saver/pkg/core"
	"password-saver/pkg/dto"
	"password-saver/pkg/logs"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepository core.UserRepository
}

func newUserUseCase(ur core.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepository: ur,
	}
}

func (uc *UserUseCase) Registration(req *dto.AuthRequest) (int64, error) {

	if err := validateAuthRequest(req); err != nil {
		logrus.Errorf(logs.FailedToValidateUser, err)
		return 0, ErrInvalidInput
	}

	hashPassword, err := hashPassword(req.Password)
	if err != nil {
		logrus.Errorf(logs.FailedToHashPassword, err)
		return 0, ErrHashPassword
	}

	regDate := getTodayDate()
	user := newUser(0, req.Email, hashPassword, regDate)

	userID, err := uc.UserRepository.Registration(user)
	if err != nil {
		return 0, handleRepositoryErrors(err)
	}

	return userID, nil
}

func (uc *UserUseCase) LogIn(req *dto.AuthRequest) (*core.User, error) {

	if err := validateAuthRequest(req); err != nil {
		logrus.Errorf(logs.FailedToValidateUser, err)
		return nil, ErrInvalidInput
	}

	user, err := uc.UserRepository.LogIn(req)
	if err != nil {
		return nil, handleRepositoryErrors(err)
	}

	if !comparePassword(req.Password, user.HashPassword) {
		logrus.Errorf(logs.FailedToComparePasswords, err)
		return nil, ErrComparePasswords
	}

	sanitizeUserStruct(user)
	return user, nil
}

func (uc *UserUseCase) Update(req *dto.UpdateUserRequest, userID int64) error {

	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		return handleRepositoryErrors(err)
	}

	if err := validateUpdateRequest(req); err != nil {
		logrus.Errorf(logs.FailedToValidateUser, err)
		return ErrInvalidInput
	}

	if !comparePassword(req.OldPassword, user.HashPassword) {
		logrus.Errorf(logs.FailedToComparePasswords, err)
		return ErrComparePasswords
	}

	user.HashPassword, err = hashPassword(req.NewPassword)
	if err != nil {
		logrus.Errorf(logs.FailedToHashPassword, err)
		return ErrHashPassword
	}

	if err := uc.UserRepository.Update(user); err != nil {
		return handleRepositoryErrors(err)
	}

	return nil
}

func (uc *UserUseCase) GetByID(userID int64) (*core.User, error) {
	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		return nil, handleRepositoryErrors(err)
	}
	sanitizeUserStruct(user)

	return user, nil
}

func (uc *UserUseCase) Delete(userID int64) error {
	if err := uc.UserRepository.Delete(userID); err != nil {
		return handleRepositoryErrors(err)
	}

	return nil
}

func validateAuthRequest(req *dto.AuthRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func validateUpdateRequest(req *dto.UpdateUserRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
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

func sanitizeUserStruct(u *core.User) {
	u.HashPassword = ""
}

func newUser(userID int64, email, hashPassword, regDate string) *core.User {
	return &core.User{
		UserID:       userID,
		Email:        email,
		HashPassword: hashPassword,
		RegDate:      regDate,
	}
}
