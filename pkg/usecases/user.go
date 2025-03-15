package usecases

import (
	"password-saver/pkg/dto"
	apperrors "password-saver/pkg/errors"
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
		logrus.Errorf("failed to validate user: %v", err)
		return 0, apperrors.ErrValidateUser
	}

	hashPassword, err := hashPassword(req.Password)
	if err != nil {
		logrus.Errorf("failed to hash password: %v", err)
		return 0, apperrors.ErrHashPassword
	}

	regDate := getTodayDate()

	user := newUser(0, req.Email, hashPassword, regDate)

	userID, err := uc.UserRepository.Registration(user)
	if err != nil {
		return 0, handleRepositoryError(err, req.Email)
	}

	logrus.Infof("user was registated successfully with id = %d", userID)

	return userID, nil
}

func (uc *UserUseCase) LogIn(req *dto.AuthRequest) (*model.User, error) {

	if err := validateAuthRequest(req); err != nil {
		logrus.Errorf("failed to validate user: %v", err)
		return nil, apperrors.ErrValidateUser
	}

	user, err := uc.UserRepository.LogIn(req)
	if err != nil {
		return nil, handleRepositoryError(err, req.Email)
	}

	if !comparePassword(req.Password, user.HashPassword) {
		logrus.Errorf("failed to compare passwords: %v", err)
		return nil, apperrors.ErrComparePasswords
	}

	sanitizeUserStruct(user)

	logrus.Infof("user {id = %d} was lodin successfully", user.UserID)

	return user, nil
}

func (uc *UserUseCase) Update(req *dto.UpdateUserRequest, userID int64) error {

	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		return handleRepositoryError(err, userID)
	}

	if err := validateUpdateRequest(req); err != nil {
		logrus.Errorf("failed to validate user: %v", err)
		return apperrors.ErrValidateUser
	}

	if !comparePassword(req.OldPassword, user.HashPassword) {
		logrus.Errorf("failed to compare passwords: %v", err)
		return apperrors.ErrComparePasswords
	}

	user.HashPassword, err = hashPassword(req.NewPassword)
	if err != nil {
		logrus.Errorf("failed to hash password: %v", err)
		return apperrors.ErrHashPassword
	}

	if err := uc.UserRepository.Update(user); err != nil {
		return handleRepositoryError(err, userID)
	}

	logrus.Infof("user {id = %d} was updated successfully", user.UserID)

	return nil
}

func (uc *UserUseCase) GetByID(userID int64) (*model.User, error) {
	user, err := uc.UserRepository.GetByID(userID)
	if err != nil {
		return nil, handleRepositoryError(err, userID)
	}
	sanitizeUserStruct(user)

	logrus.Info("successfull getting by id")
	return user, nil
}

func (uc *UserUseCase) Delete(userID int64) error {
	if err := uc.UserRepository.Delete(userID); err != nil {
		return handleRepositoryError(err, userID)
	}

	logrus.Infof("user {id = %d} was deleted successfully", userID)
	return nil
}

func validateAuthRequest(req *dto.AuthRequest) error {
	validate := validator.New()
	err := validate.Struct(req)
	return err
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

func handleRepositoryError(err error, data interface{}) error {
	switch err {
	case apperrors.ErrUserNotFound:
		logrus.Errorf("userID: %d %v", data, err)
		return apperrors.ErrUserNotFound

	case apperrors.ErrDuplicateUser:
		logrus.Errorf("email: %s %v", data, err)
		return apperrors.ErrDuplicateUser

	default:
		logrus.Errorf("internal database error: %v", err)
		return apperrors.ErrDatabaseInternal
	}
}
