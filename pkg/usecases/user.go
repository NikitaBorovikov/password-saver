package usecases

import (
	"errors"
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	errComparePasswords = errors.New("failed compare passwords: incorrected password")
	errValidateUser     = errors.New("failed to validate user")
	errHashPassword     = errors.New("failed to hash password")
)

type UserUseCase struct {
	UserRepository model.UserRepository
}

func NewUserUseCase(ur model.UserRepository) *UserUseCase {
	return &UserUseCase{
		UserRepository: ur,
	}
}

func (uc *UserUseCase) Registration(req *dto.RegRequest) (int64, error) {

	if err := validateRegRequest(req); err != nil {
		logError(req.Email, err, "failed to validate user")
		return 0, err
	}

	hashPassword, err := hashPassword(req.Password)
	if err != nil {
		logError(req.Email, err, "failed to hash password")
		return 0, err
	}

	user := &model.User{
		Email:        req.Email,
		HashPassword: hashPassword,
		RegDate:      getTodayDate(),
	}

	userID, err := uc.UserRepository.Registration(user)
	if err != nil {
		logError(req.Email, err, "failed to save user")
		return 0, err
	}

	logInfo(userID, "user was registered successfully")

	return userID, nil

}

func (uc *UserUseCase) LogIn(req *dto.LogInRequest) (*model.User, error) {

	if err := validateLoginRequest(req); err != nil {
		logError(req.Email, err, "failed to validate user")
		return nil, err
	}

	user, err := uc.UserRepository.LogIn(req)
	if err != nil {
		logError(req.Email, err, "failed to login user")
		return nil, err
	}

	if !comparePassword(req.Password, user.HashPassword) {
		logError(req.Email, errComparePasswords, "failed compare passwords")
		return nil, errComparePasswords
	}

	sanitizeUserStruct(user)

	logInfo(user.UserID, "successful login")

	return user, nil
}

func (uc *UserUseCase) Update(req *dto.UpdateUserRequest, userID int64) error {

	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		logErrorWithID(userID, err, "failed to get user by ID")
		return err
	}

	if err := validateUpdateRequest(req); err != nil {
		logErrorWithID(userID, err, "failed to validate user")
		return err
	}

	if !comparePassword(req.OldPassword, user.HashPassword) {
		logErrorWithID(userID, err, "failed to compare passwords")
		return errComparePasswords
	}

	user.HashPassword, err = hashPassword(req.NewPassword)
	if err != nil {
		logErrorWithID(userID, err, "failed to hash password")
		return err
	}

	if err := uc.UserRepository.Update(user); err != nil {
		logErrorWithID(userID, err, "failed to update user")
		return err
	}

	logInfo(user.UserID, "user was updated successfully")

	return nil
}

func (uc *UserUseCase) GetByID(userID int64) (*model.User, error) {
	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		logErrorWithID(userID, err, "failed to get user by ID")
		return nil, err
	}
	sanitizeUserStruct(user)
	logInfo(userID, "successful get by ID")
	return user, nil
}

func (uc *UserUseCase) Delete(userID int64) error {
	if err := uc.UserRepository.Delete(userID); err != nil {
		logErrorWithID(userID, err, "failed to delete user")
		return err
	}

	logInfo(userID, "user was deleted successfully")
	return nil
}

func validateRegRequest(req *dto.RegRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}

	return nil
}

func validateLoginRequest(req *dto.LogInRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return errValidateUser
	}

	return nil
}

func validateUpdateRequest(req *dto.UpdateUserRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return errValidateUser
	}

	return nil
}

func getTodayDate() string {
	return time.Now().Format(time.RFC3339)
}

func hashPassword(inputPassword string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", errHashPassword
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

func logInfo(userID int64, msg string) {
	logrus.WithFields(logrus.Fields{
		"userID": userID,
	}).Info(msg)
}

func logError(email string, err error, msg string) {
	logrus.WithFields(logrus.Fields{
		"email": email,
		"error": err,
	}).Error(msg)
}

func logErrorWithID(userID int64, err error, msg string) {
	logrus.WithFields(logrus.Fields{
		"email": userID,
		"error": err,
	}).Error(msg)
}
