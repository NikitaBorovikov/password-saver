package usecases

import (
	"crypto/rand"
	"encoding/hex"
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

	salt, err := generateUserSalt()
	if err != nil {
		return 0, err
	}

	hashPassword, err := hashPassword(req.Password, salt)
	if err != nil {
		return 0, err
	}

	user := &model.User{
		Email:        req.Email,
		HashPassword: hashPassword,
		Salt:         salt,
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

	if !comparePassword(q.Password, user) {
		return nil, fmt.Errorf("failed compare passwords: incorrected password")
	}

	sanitizeUserStruct(user)

	return user, nil
}

func (uc *UserUseCase) Update(req *dto.UpdateUserRequest) error {

	user, err := uc.UserRepository.GetUserByID(req.UserID)
	if err != nil {
		return err
	}

	if err := validateForUpdateUser(req); err != nil {
		return err
	}

	if !comparePassword(req.OldPassword, user) {
		return fmt.Errorf("failed compare passwords: incorrected password")
	}

	salt, err := generateUserSalt()
	if err != nil {
		return err
	}

	hashPassword, err := hashPassword(req.NewPassword, salt)
	if err != nil {
		return err
	}

	user.Salt = salt
	user.HashPassword = hashPassword

	if err := uc.UserRepository.Update(user); err != nil {
		return err
	}

	return nil
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

func generateUserSalt() (string, error) {
	byteArr := make([]byte, 32)

	_, err := rand.Read(byteArr)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %v", err)
	}

	salt := hex.EncodeToString(byteArr)
	return salt, nil
}

func hashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash of the password: %v", err)
	}

	return string(hashPassword), nil
}

func comparePassword(inputPassword string, u *model.User) bool {
	inputPasswordWithSalt := inputPassword + u.Salt
	if err := bcrypt.CompareHashAndPassword([]byte(u.HashPassword), []byte(inputPasswordWithSalt)); err != nil {
		return false
	}
	return true
}

func sanitizeUserStruct(u *model.User) {
	u.HashPassword = ""
	u.Salt = ""
}

func (uc *UserUseCase) GetUserByID(userID int64) (*model.User, error) {
	user, err := uc.UserRepository.GetUserByID(userID)
	return user, err
}
