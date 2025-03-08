package usecases

import (
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

func (uc *PasswordUseCase) Save(req *dto.SavePasswordRequest) (err error) {

	if err := validateForSavePassword(req); err != nil {
		return err
	}

	var p model.Password

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

func (uc *PasswordUseCase) GetAll(userID int64) ([]model.Password, error) {
	return nil, nil
}

func (uc *PasswordUseCase) GetByID(passwordID string) (*model.Password, error) {
	return nil, nil
}

func (uc *PasswordUseCase) Update(p *model.Password) error {
	return nil
}

func (uc *PasswordUseCase) Delete(passwordID string) error {
	return nil
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
