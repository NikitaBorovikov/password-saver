package usecases

import (
	"encoding/base64"
	"fmt"
	"password-saver/pkg/config"
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
	"password-saver/pkg/usecases/encryption"

	"github.com/go-playground/validator"
)

type PasswordUseCase struct {
	PasswordRepository model.PasswordRepository
	cfg                *config.EncryptKeys
}

func NewPasswordUseCase(pr model.PasswordRepository, cfg *config.EncryptKeys) *PasswordUseCase {
	return &PasswordUseCase{
		PasswordRepository: pr,
		cfg:                cfg,
	}
}

func (uc *PasswordUseCase) Save(req *dto.PasswordRequest, userID int64) error {

	if err := validateForSavePassword(req); err != nil {
		return err
	}

	encPassword, encService, err := uc.encryptFields(req)
	if err != nil {
		return err
	}

	password := newPassword(0, userID, encPassword, encService)

	if err := uc.PasswordRepository.Save(password); err != nil {
		return err
	}
	return nil
}

func (uc *PasswordUseCase) GetAll(userID int64) ([]dto.PasswordResponse, error) {
	userPasswords, err := uc.PasswordRepository.GetAll(userID)
	if err != nil {
		return nil, err
	}

	passwordResponse, err := uc.makePasswordResponse(userPasswords)
	if err != nil {
		return nil, err
	}

	return passwordResponse, nil
}

func (uc *PasswordUseCase) GetByID(passwordID int64) (*dto.PasswordResponse, error) {
	userPassword, err := uc.PasswordRepository.GetByID(passwordID)
	if err != nil {
		return nil, err
	}

	passwordResponse, err := uc.decryptFields(*userPassword)
	if err != nil {
		return nil, err
	}

	return passwordResponse, nil
}

func (uc *PasswordUseCase) Update(req *dto.PasswordRequest, passwordID, userID int64) error {

	if err := validateForUpdatePassword(req); err != nil {
		return err
	}

	encPassword, encService, err := uc.encryptFields(req)
	if err != nil {
		return err
	}

	password := newPassword(passwordID, userID, encPassword, encService)

	if err := uc.PasswordRepository.Update(password); err != nil {
		return err
	}

	return nil
}

func (uc *PasswordUseCase) Delete(passwordID int64) error {
	err := uc.PasswordRepository.Delete(passwordID)
	return err
}

func (uc *PasswordUseCase) makePasswordResponse(userPasswords []model.Password) ([]dto.PasswordResponse, error) {
	passwordResponse := make([]dto.PasswordResponse, 0, len(userPasswords))

	for _, elem := range userPasswords {

		userPassword, err := uc.decryptFields(elem)
		if err != nil {
			return nil, err
		}

		passwordResponse = append(passwordResponse, *userPassword)
	}

	return passwordResponse, nil
}

func (uc *PasswordUseCase) encryptFields(req *dto.PasswordRequest) (string, string, error) {

	encPassword, err := encryption.Encrypt([]byte(req.Password), []byte(uc.cfg.EncPasswordKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt password: %v", err)
	}

	encService, err := encryption.Encrypt([]byte(req.Service), []byte(uc.cfg.EncServiceKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt service: %v", err)
	}

	return encPassword, encService, nil
}

func (uc *PasswordUseCase) decryptFields(password model.Password) (*dto.PasswordResponse, error) {

	var passwordResponse dto.PasswordResponse
	var err error

	passwordResponse.Password, err = decryptData(password.EncPassword, uc.cfg.EncPasswordKey)
	if err != nil {
		return nil, err
	}

	passwordResponse.Service, err = decryptData(password.EncService, uc.cfg.EncServiceKey)
	if err != nil {
		return nil, err
	}

	return &passwordResponse, nil
}

func validateForSavePassword(req *dto.PasswordRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("failed to validate password struct: %v", err)
	}
	return nil
}

func validateForUpdatePassword(req *dto.PasswordRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("failed to validate password struct: %v", err)
	}
	return nil
}

func decryptData(encData string, encKey string) (string, error) {
	encDataInByte, err := base64.StdEncoding.DecodeString(encData)
	if err != nil {
		return "", fmt.Errorf("failed to decode string by byte: %v", err)
	}

	plainData, err := encryption.Decrypt(encDataInByte, []byte(encKey))
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %v", err)
	}

	return plainData, nil
}

func newPassword(passwordID, userID int64, encPassword, encService string) *model.Password {
	return &model.Password{
		PasswordID:  passwordID,
		UserID:      userID,
		EncPassword: encPassword,
		EncService:  encService,
	}
}
