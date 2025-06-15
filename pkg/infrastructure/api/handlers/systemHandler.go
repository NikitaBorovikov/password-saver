package handlers

import (
	"net/http"
	"password-saver/pkg/usecases"

	"github.com/sirupsen/logrus"
)

type SystemHandler struct {
	SystemUseCase *usecases.SystemUseCase
}

func newSystemHandler(uc *usecases.SystemUseCase) *SystemHandler {
	return &SystemHandler{
		SystemUseCase: uc,
	}
}

// @Summary Health Checking
// @Description Сhecks the API operation status.
// @Tags System
// @Produce json
// @Success 200 {object} dto.HealthCheckResponse
// @Failure 503 {object} dto.HealthCheckResponse
// @Router /health [get]
func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	healthCheckResponse, err := h.SystemUseCase.HealthCheck()
	if err != nil {
		sendOKResponse(w, r, http.StatusServiceUnavailable, healthCheckResponse)
		logrus.Errorf("health checking with err: %v", err)
		return
	}

	sendOKResponse(w, r, http.StatusOK, healthCheckResponse)
	logrus.Info("health checking: ok")
}
