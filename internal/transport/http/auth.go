package httptransport

import (
	"net/http"

	commondomain "github.com/SashaMaltsev/room-booking-service/internal/domain/common"
	"github.com/SashaMaltsev/room-booking-service/internal/transport/http/dto"
)

func (a *API) handleRoot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) handleInfo(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) handleDummyLogin(w http.ResponseWriter, r *http.Request) {
	if a.auth == nil {
		writeDomainError(w, nil)
		return
	}

	var request dto.DummyLoginRequest
	if err := decodeJSON(r, &request); err != nil {
		writeInvalidRequest(w, "invalid request body")
		return
	}

	role, err := commondomain.ParseRole(request.Role)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	token, err := a.auth.DummyLogin(r.Context(), role)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.TokenResponse{Token: token})
}
