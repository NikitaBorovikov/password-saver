package usecases

import (
	"fmt"
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
	"time"

	"github.com/go-playground/validator"
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

func (uc *UserUseCase) Registration(req *dto.RegRequest) (int64, error) {

	if err := validateForRegistration(req); err != nil {
		return 0, err
	}

	hashPassword, err := hashPassword(req.Password)
	if err != nil {
		return 0, err
	}

	user := &model.User{
		Email:        req.Email,
		HashPassword: hashPassword,
		RegDate:      getTodayDate(),
	}

	userID, err := uc.UserRepository.Registration(user)
	return userID, err

}

func (uc *UserUseCase) LogIn(q *dto.LogInRequest) (*model.User, error) {

	user, err := uc.UserRepository.LogIn(q)
	if err != nil {
		return nil, err
	}

	if !comparePassword(q.Password, user.HashPassword) {
		return nil, fmt.Errorf("failed compare passwords: incorrected password")
	}

	sanitizeUserStruct(user)

	return user, nil
}

func (uc *UserUseCase) Update(req *dto.UpdateUserRequest) error {

	user, err := uc.UserRepository.GetByID(req.UserID)
	if err != nil {
		return err
	}

	if err := validateForUpdateUser(req); err != nil {
		return err
	}

	if !comparePassword(req.OldPassword, user.HashPassword) {
		return fmt.Errorf("failed compare passwords: incorrected password")
	}

	user.HashPassword, err = hashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := uc.UserRepository.Update(user); err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) GetByID(userID int64) (*model.User, error) {
	user, err := uc.UserRepository.GetByID(userID)
	user.HashPassword = ""
	return user, err
}

func (uc *UserUseCase) Delete(userID int64) error {
	err := uc.UserRepository.Delete(userID)
	return err
}

func validateForRegistration(req *dto.RegRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("failed to validate user struct: %v", err)
	}

	return nil
}

func validateForUpdateUser(req *dto.UpdateUserRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("failed to validate user struct: %v", err)
	}

	return nil
}

func getTodayDate() string {
	return time.Now().Format(time.RFC3339)
}

func hashPassword(inputPassword string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash of the password: %v", err)
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
