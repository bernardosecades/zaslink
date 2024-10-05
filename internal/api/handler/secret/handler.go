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
	ctx := r.Context()
	var request CreateSecretRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		api.EncodeHTTPError(api.NewHTTPError("wrong request parameters, please check the body request", http.StatusBadRequest), w)
		return
	}

	secret, err := h.secretService.CreateSecret(ctx, request.Content, request.Pwd, request.Expiration)
	if err != nil {
		api.EncodeHTTPError(api.NewHTTPError(err.Error(), http.StatusBadRequest), w)
		return
	}

	response := CreateSecretResponse{
		ID:        secret.ID,
		PrivateID: secret.PrivateID,
		ExpiredAt: secret.ExpiredAt,
	}

	api.EncodeHTTPResponse(response, w, http.StatusCreated)
}

func (h *Handler) RetrieveSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pwd := r.Header.Get("X-Password")
	ID := vars["id"]
	ctx := r.Context()

	secret, err := h.secretService.RetrieveSecret(ctx, ID, pwd)
	if err != nil {
		if errors.Is(err, service.ErrSecretDoesNotExist) {
			api.EncodeHTTPError(api.NewHTTPError("The secret does not exist or has already been read.", http.StatusNotFound), w)
			return
		}
		if errors.Is(err, service.ErrInvalidPassword) {
			api.EncodeHTTPError(api.NewHTTPError(service.ErrInvalidPassword.Error(), http.StatusBadRequest), w)
			return
		}
		api.EncodeHTTPError(api.NewHTTPError("There was an error reading the secret. Please try again later.", http.StatusInternalServerError), w)
		return
	}

	response := RetrieveSecretResponse{
		Content: string(secret.Content),
	}

	api.EncodeHTTPResponse(response, w, http.StatusOK)
}

func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	privateID := vars["private_id"]
	ctx := r.Context()

	err := h.secretService.DeleteSecret(ctx, privateID)
	if err != nil {
		if errors.Is(err, service.ErrSecretDoesNotExist) {
			api.EncodeHTTPError(api.NewHTTPError("The secret does not exist ", http.StatusNotFound), w)
			return
		}
		api.EncodeHTTPError(api.NewHTTPError("There was an error reading the secret. Please try again later.", http.StatusInternalServerError), w)
		return
	}

	api.EncodeHTTPResponse(nil, w, http.StatusNoContent)
}
