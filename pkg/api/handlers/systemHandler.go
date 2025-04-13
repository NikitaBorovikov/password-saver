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

func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.SystemUseCase.PingDb(); err != nil {
		logrus.Errorf("health check: ping db err: %v", err)
		sendErrorRespose(w, r, http.StatusServiceUnavailable, err)
		return
	}
	sendOKResponse(w, r, http.StatusOK, nil)
	logrus.Info("health check: ok")
}
