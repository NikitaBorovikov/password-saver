package usecases

import (
	"encoding/base64"
	"fmt"
	"password-saver/pkg/config"
	"password-saver/pkg/dto"
	"password-saver/pkg/logs"
	"password-saver/pkg/model"
	"password-saver/pkg/usecases/encryption"
	"password-saver/pkg/usecases/generation"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
)

type PasswordUseCase struct {
	PasswordRepository model.PasswordRepository
	cfg                *config.EncryptKeys
}

func newPasswordUseCase(pr model.PasswordRepository, cfg *config.EncryptKeys) *PasswordUseCase {
	return &PasswordUseCase{
		PasswordRepository: pr,
		cfg:                cfg,
	}
}

type encPasswordData struct {
	password string
	service  string
	login    string
}

func (uc *PasswordUseCase) Save(req *dto.PasswordRequest, userID int64) error {

	if err := validateForPassword(req); err != nil {
		logrus.Errorf(logs.FailedToValidatePassword, err)
		return ErrInvalidInput
	}

	encPasswordData, err := uc.encryptFields(req)
	if err != nil {
		logrus.Error(err)
		return ErrEcryptData
	}

	password := newPassword(0, userID, encPasswordData)

	if err := uc.PasswordRepository.Save(password); err != nil {
		return handleRepositoryErrors(err)
	}
	return nil
}

func (uc *PasswordUseCase) GetAll(userID int64) ([]dto.PasswordResponse, error) {
	userPasswords, err := uc.PasswordRepository.GetAll(userID)
	if err != nil {
		return nil, handleRepositoryErrors(err)
	}

	passwordResponse, err := uc.makePasswordResponse(userPasswords)
	if err != nil {
		return nil, ErrMakePasswordResponse
	}

	return passwordResponse, nil
}

func (uc *PasswordUseCase) GetByID(passwordID, userID int64) (*dto.PasswordResponse, error) {
	userPassword, err := uc.PasswordRepository.GetByID(passwordID, userID)
	if err != nil {
		return nil, handleRepositoryErrors(err)
	}

	passwordResponse, err := uc.decryptFields(*userPassword)
	if err != nil {
		return nil, ErrDecryptData
	}

	passwordResponse.PasswordID = passwordID

	return passwordResponse, nil
}

func (uc *PasswordUseCase) Update(req *dto.PasswordRequest, passwordID, userID int64) error {

	if err := validateForPassword(req); err != nil {
		logrus.Errorf(logs.FailedToValidatePassword, err)
		return ErrInvalidInput
	}

	encPasswordData, err := uc.encryptFields(req)
	if err != nil {
		logrus.Error(err)
		return ErrEcryptData
	}

	password := newPassword(0, userID, encPasswordData)

	if err := uc.PasswordRepository.Update(password); err != nil {
		return handleRepositoryErrors(err)
	}

	return nil
}

func (uc *PasswordUseCase) Delete(passwordID, userID int64) error {
	if err := uc.PasswordRepository.Delete(passwordID, userID); err != nil {
		return handleRepositoryErrors(err)
	}
	return nil
}

func (uc *PasswordUseCase) Generate(ps *dto.GeneratePasswordRequest) (string, error) {

	if err := validateGenPasswordSettings(ps); err != nil {
		logrus.Errorf(logs.FailedToValidatePasswordSettings, err)
		return "", ErrInvalidInput
	}

	password := generation.GenNewPassword(ps)

	return password, nil
}

func (uc *PasswordUseCase) makePasswordResponse(userPasswords []model.Password) ([]dto.PasswordResponse, error) {
	passwordResponse := make([]dto.PasswordResponse, 0, len(userPasswords))

	for _, elem := range userPasswords {

		userPassword, err := uc.decryptFields(elem)
		if err != nil {
			return nil, err
		}
		userPassword.PasswordID = elem.PasswordID

		passwordResponse = append(passwordResponse, *userPassword)
	}

	return passwordResponse, nil
}

func (uc *PasswordUseCase) encryptFields(req *dto.PasswordRequest) (*encPasswordData, error) {
	epd := &encPasswordData{}
	var err error

	epd.password, err = encryption.Encrypt([]byte(req.Password), []byte(uc.cfg.EncPasswordKey))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %v", err)
	}

	epd.service, err = encryption.Encrypt([]byte(req.Service), []byte(uc.cfg.EncServiceKey))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt service: %v", err)
	}

	epd.login, err = encryption.Encrypt([]byte(req.Login), []byte(uc.cfg.EncLoginKey))
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt login: %v", err)
	}

	return epd, nil
}

func (uc *PasswordUseCase) decryptFields(password model.Password) (*dto.PasswordResponse, error) {

	var passwordResponse dto.PasswordResponse
	var err error

	passwordResponse.Password, err = decryptData(password.EncPassword, uc.cfg.EncPasswordKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %v", err)
	}

	passwordResponse.Service, err = decryptData(password.EncService, uc.cfg.EncServiceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt service: %v", err)
	}

	// if login is not NULL, then we decrypt it, otherwise we set the default value
	if password.EncLogin != nil {
		passwordResponse.Login, err = decryptData(*password.EncLogin, uc.cfg.EncLoginKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt login: %v", err)
		}
	} else {
		passwordResponse.Login = ""
	}

	return &passwordResponse, nil
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

func validateForPassword(req *dto.PasswordRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func validateGenPasswordSettings(req *dto.GeneratePasswordRequest) error {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func newPassword(passwordID, userID int64, ecp *encPasswordData) *model.Password {
	return &model.Password{
		PasswordID:  passwordID,
		UserID:      userID,
		EncPassword: ecp.password,
		EncService:  ecp.service,
		EncLogin:    &ecp.login,
	}
}
