package usecases

import "password-saver/pkg/model"

type SystemUseCase struct {
	SystemRepository model.SystemRepository
}

func newSystemUseCase(sr model.SystemRepository) *SystemUseCase {
	return &SystemUseCase{
		SystemRepository: sr,
	}
}

func (su *SystemUseCase) PingDb() error {
	return su.SystemRepository.PingDB()
}
