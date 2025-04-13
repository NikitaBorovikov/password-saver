package usecases

import (
	"password-saver/pkg/dto"
	"password-saver/pkg/model"
	"time"
)

type SystemUseCase struct {
	SystemRepository model.SystemRepository
}

func newSystemUseCase(sr model.SystemRepository) *SystemUseCase {
	return &SystemUseCase{
		SystemRepository: sr,
	}
}

func (su *SystemUseCase) HealthCheck() (*dto.HealthCheckResponse, error) {

	resp := dto.HealthCheckResponse{}
	resp.Time = time.Now()

	if err := su.SystemRepository.PingDB(); err != nil {
		resp.Status = "unhealthy"
		resp.Details = "databse: failed to ping"
		return &resp, err
	}

	resp.Status = "healthy"
	return &resp, nil
}
