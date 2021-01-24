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
	err := json.NewDecoder(r.Body).Decode(&cs)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	pass := r.Header.Get("X-Password")
	if pass == "" {
		pass = os.Getenv("SECRET_PASSWORD")
	}

	secret, err := controller.secretService.CreateSecret(cs.Content, pass)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	// TODO change url
	cr := dto.CreateSecretResponse{
		Url: "http://127.0.0.1:8080/secret/" + secret.Id,
	}

	err = json.NewEncoder(w).Encode(cr)
	if err != nil {
		panic("error to encode create secret response")
	}
}

func (controller *SecretController) GetSecret(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pass := r.Header.Get("X-Password")

	if pass == "" {
		pass = os.Getenv("SECRET_PASSWORD")
	}

	content, err := controller.secretService.GetContentSecret(id, pass)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	sr := dto.SecretResponse{content}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(sr)
	if err != nil {
		panic("error to encode secret")
	}
}
