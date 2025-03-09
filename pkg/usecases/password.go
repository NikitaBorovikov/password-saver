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

func (uc *PasswordUseCase) Save(req *dto.SavePasswordRequest, userID int64) (err error) {

	if err := validateForSavePassword(req); err != nil {
		return err
	}

	var p model.Password
	p.UserID = userID

	p.EncPassword, err = encryptData(req.Password, uc.cfg.EncPasswordKey)
	if err != nil {
		return err
	}

	p.EncService, err = encryptData(req.Service, uc.cfg.EncServiceKey)
	if err != nil {
		return err
	}

	if err := uc.PasswordRepository.Save(&p); err != nil {
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

	var passwordResponse dto.PasswordResponse

	passwordResponse.Password, err = decodeData(userPassword.EncPassword, uc.cfg.EncPasswordKey)
	if err != nil {
		return nil, err
	}

	passwordResponse.Service, err = decodeData(userPassword.EncService, uc.cfg.EncServiceKey)
	if err != nil {
		return nil, err
	}

	return &passwordResponse, nil
}

func (uc *PasswordUseCase) Update(p *model.Password) error {
	return nil
}

func (uc *PasswordUseCase) Delete(passwordID int64) error {
	err := uc.PasswordRepository.Delete(passwordID)
	return err
}

func (uc *PasswordUseCase) makePasswordResponse(userPasswords []model.Password) ([]dto.PasswordResponse, error) {
	passwordResponse := make([]dto.PasswordResponse, 0, len(userPasswords))
	var err error

	for _, elem := range userPasswords {

		var userPassword dto.PasswordResponse

		userPassword.Password, err = decodeData(elem.EncPassword, uc.cfg.EncPasswordKey)
		if err != nil {
			return nil, err
		}

		userPassword.Service, err = decodeData(elem.EncService, uc.cfg.EncServiceKey)
		if err != nil {
			return nil, err
		}

		passwordResponse = append(passwordResponse, userPassword)
	}

	return passwordResponse, nil
}

func validateForSavePassword(req *dto.SavePasswordRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return fmt.Errorf("failed to validate password struct: %v", err)
	}
	return nil
}

func encryptData(data string, encKey string) (string, error) {
	encData, err := encryption.Encrypt([]byte(data), []byte(encKey))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %v", err)
	}

	return encData, nil
}

func decodeData(encData string, encKey string) (string, error) {
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
