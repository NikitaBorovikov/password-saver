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

func (uc *UserUseCase) Registration(u *model.User) error {

	setRegDateForUser(u)

	if err := validateForRegistration(u); err != nil {
		return err
	}

	if err := generateUserSalt(u); err != nil {
		return err
	}

	if err := hashPassword(u); err != nil {
		return err
	}

	if err := uc.UserRepository.Registration(u); err != nil {
		return err
	}

	return nil
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

func (uc *UserUseCase) Update(u *model.User) error {
	return nil
}

func (uc *UserUseCase) Delete(userID int64) error {
	return nil
}

func validateForRegistration(u *model.User) error {
	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		return fmt.Errorf("failed to validate user struct: %v", err)
	}

	return nil
}

func setRegDateForUser(u *model.User) {
	u.RegDate = time.Now().Format(time.RFC3339)
}

func generateUserSalt(u *model.User) error {
	byteArr := make([]byte, 32)

	_, err := rand.Read(byteArr)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}

	u.Salt = hex.EncodeToString(byteArr)
	return nil
}

func hashPassword(u *model.User) error {
	saltedPassword := u.Password + u.Salt

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash of the password: %v", err)
	}

	u.HashPassword = string(hashPassword)
	return nil
}

func comparePassword(inputPassword, hashPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(inputPassword)); err != nil {
		return false
	}
	return true
}

func sanitizeUserStruct(u *model.User) {
	u.Password = ""
	u.HashPassword = ""
	u.Salt = ""
}
