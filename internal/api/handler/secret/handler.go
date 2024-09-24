package secret

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bernardosecades/sharesecret/pkg/api"

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
	var request CreateSecretRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.EncodeHTTPError(api.NewHTTPError("wrong request parameters, please check the body request", http.StatusBadRequest), w)
		return
	}

	secret, err := h.secretService.CreateSecret(r.Context(), request.Content, request.Pwd)
	if err != nil {
		api.EncodeHTTPError(api.NewHTTPError(err.Error(), http.StatusBadRequest), w)
		return
	}

	response := CreateSecretResponse{
		ID:        secret.ID,
		ExpiredAt: secret.ExpiredAt,
	}

	api.EncodeHTTPResponse(response, w, http.StatusCreated)
}

func (h *Handler) RetrieveSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pwd := r.Header.Get("X-Password")
	ID := vars["id"]

	secret, err := h.secretService.RetrieveSecret(r.Context(), ID, pwd)
	if err != nil {
		if errors.Is(err, service.ErrSecretDoesNotExist) {
			api.EncodeHTTPError(api.NewHTTPError("secret does not exist or already read", http.StatusNotFound), w)
			return
		}
		api.EncodeHTTPError(api.NewHTTPError("there was an error to read the secret, try later", http.StatusInternalServerError), w)
		return
	}

	response := RetrieveSecretResponse{
		Content: secret.Content,
	}

	api.EncodeHTTPResponse(response, w, http.StatusOK)
}
