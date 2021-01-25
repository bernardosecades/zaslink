package controllers

import (
	"encoding/json"
	"github.com/bernardosecades/sharesecret/dto"
	"github.com/bernardosecades/sharesecret/services"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

type SecretController struct {
	secretService services.SecretService
}

func NewSecretController(s services.SecretService) *SecretController {
	return &SecretController{s}
}

func (controller *SecretController) CreateSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var cs dto.CreateSecretRequest
	_ = json.NewDecoder(r.Body).Decode(&cs)

	if len(cs.Content) == 0 {
		customError := dto.ErrorResponse{StatusCode: 400, Err: "Empty 'content'"}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	customPwd := true
	pass := r.Header.Get("X-Password")
	if pass == "" {
		pass = os.Getenv("SECRET_DEFAULT_PASSWORD")
		customPwd = false
	}

	// TODO default password injected in sercice so refactor (pass will be optional will simplify the controller) if X-Password is empty ignore him
	secret, err := controller.secretService.CreateSecret(cs.Content, pass, customPwd)

	if err != nil {
		customError := dto.ErrorResponse{StatusCode: 404, Err: err.Error()}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	cr := dto.CreateSecretResponse{
		Url: os.Getenv("SERVER_URL") + ":" + os.Getenv("SERVER_PORT") + "/secret/" + secret.Id,
	}
	_ = json.NewEncoder(w).Encode(cr)
}

func (controller *SecretController) GetSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	pass := r.Header.Get("X-Password")

	customPwd, err := controller.secretService.HasSecretWithCustomPwd(id)

	if err != nil {
		customError := dto.ErrorResponse{StatusCode: 404, Err: err.Error()}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	if customPwd && pass == "" {
		customError := dto.ErrorResponse{StatusCode: 400, Err: "Missing X-Password"}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	if pass == "" {
		pass = os.Getenv("SECRET_DEFAULT_PASSWORD")
	}

	content, err := controller.secretService.GetContentSecret(id, pass)

	if err != nil {
		customError := dto.ErrorResponse{StatusCode: 404, Err: err.Error()}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	if len(content) == 0 {
		customError := dto.ErrorResponse{StatusCode: 400, Err: "We can not decrypt the content"}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	sr := dto.SecretResponse{content}
	_ = json.NewEncoder(w).Encode(sr)
}
