package secret

import (
	"encoding/json"
	"net/http"

	"github.com/bernardosecades/sharesecret/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	secretService *service.SecretService
}

func NewHandler(secretService *service.SecretService) *Handler {
	return &Handler{secretService: secretService}
}

func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request CreateSecretRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	secret, err := h.secretService.CreateSecret(r.Context(), request.Content, request.Pwd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	response := CreateSecretResponse{
		ID:        secret.ID,
		ExpiredAt: secret.ExpiredAt,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) RetrieveSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	pwd := r.Header.Get("X-Password")
	ID := vars["id"]

	secret, err := h.secretService.RetrieveSecret(r.Context(), ID, pwd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response := RetrieveSecretResponse{
		Content: secret.Content,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
