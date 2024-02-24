package http

import (
	"github.com/bernardosecades/sharesecret/service"
	"github.com/bernardosecades/sharesecret/types"

	"github.com/gorilla/mux"

	"encoding/json"
	"log"
	"net/http"
	"os"
)

type SecretHandler struct {
	secretService *service.SecretService
}

func NewSecretHandler(s *service.SecretService) *SecretHandler {
	return &SecretHandler{s}
}

// CreateSecret handler function to create a secret
// @Summary Create secret
// @Tags         Share Secret
// @Description Create secret
// @Produce  json
// @Param X-Password header string false "password"
// @Param request body types.CreateSecretRequest true "query params"
// @Success 201 {object} types.CreateSecretResponse
// @Failure 400 "invalid param"
// @Failure 500 "internal error"
// @Router /secret [post]
func (controller *SecretHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cs types.CreateSecretRequest
	err := json.NewDecoder(r.Body).Decode(&cs)
	if err != nil {
		customError := types.ErrorResponse{StatusCode: http.StatusBadRequest, Err: "Bad JSON in CreateSecretRequest'"}
		encodeCustomError(customError, w)
		return
	}

	if len(cs.Content) == 0 {
		customError := types.ErrorResponse{StatusCode: http.StatusBadRequest, Err: "Empty 'content'"}
		encodeCustomError(customError, w)
		return
	}

	if len(cs.Content) > 10000 {
		customError := types.ErrorResponse{StatusCode: http.StatusBadRequest, Err: "Text too long"}
		encodeCustomError(customError, w)
		return
	}

	pass := r.Header.Get("X-Password")
	secret, err := controller.secretService.CreateSecret(cs.Content, pass)
	if err != nil {
		customError := types.ErrorResponse{StatusCode: http.StatusBadRequest, Err: "We can not create your secret"}
		log.Println(err.Error())
		encodeCustomError(customError, w)
		return
	}

	cr := types.CreateSecretResponse{
		URL: os.Getenv("SERVER_URL") + "/secret/" + secret.ID,
		ID:  secret.ID,
	}

	if err := json.NewEncoder(w).Encode(cr); err != nil {
		log.Fatal("Error to encode CreateSecretResponse")
	}
}

// GetSecret handler function to get secret
// @Summary Get secret
// @Tags         Share Secret
// @Description Get secret
// @Produce  json
// @Param X-Password header string false "password"
// @Param id path string true "ID of the secret"
// @Success 200 {object} types.SecretResponse
// @Failure 400 "id is missing"
// @Failure 404 "secret can not found"
// @Failure 500 "internal error"
// @Router /secret/{id} [get]
func (controller *SecretHandler) GetSecret(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	pass := r.Header.Get("X-Password")

	customPwd, err := controller.secretService.HasSecretWithCustomPwd(id)
	if err != nil {
		customError := types.ErrorResponse{StatusCode: http.StatusNotFound, Err: "Unknown Secret: It either never existed or has already been viewed"}
		log.Println(err.Error())
		encodeCustomError(customError, w)
		return
	}

	if customPwd && pass == "" {
		customError := types.ErrorResponse{StatusCode: http.StatusBadRequest, Err: "Missing X-Password"}
		encodeCustomError(customError, w)
		return
	}

	content, err := controller.secretService.GetContentSecret(id, pass)
	if err != nil {
		customError := types.ErrorResponse{StatusCode: http.StatusNotFound, Err: "Unknown Secret: It either never existed or has already been viewed"}
		log.Println(err.Error())
		encodeCustomError(customError, w)
		return
	}

	if len(content) == 0 {
		customError := types.ErrorResponse{StatusCode: http.StatusBadRequest, Err: "We can not decrypt the content"}
		encodeCustomError(customError, w)
		return
	}

	sr := types.SecretResponse{Content: content}
	if err := json.NewEncoder(w).Encode(sr); err != nil {
		log.Fatal(err.Error())
	}
}

func encodeCustomError(customError types.ErrorResponse, w http.ResponseWriter) {
	w.WriteHeader(customError.StatusCode)
	if err := json.NewEncoder(w).Encode(customError); err != nil {
		log.Fatal("Error to encode customError")
	}
}
