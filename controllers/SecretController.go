package controllers

import (
	"encoding/json"
	"github.com/bernardosecades/sharesecret/services"
	"github.com/gorilla/mux"
	"net/http"
)

type SecretController struct {
	secretService services.SecretService
}

func NewSecretController(s services.SecretService) *SecretController {
	return &SecretController{s}
}

func (controller *SecretController) CreateSecret(w http.ResponseWriter, r *http.Request) {

	secret := controller.secretService.CreateSecret("raw content", "myPassword")

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(secret)
	if err != nil {
		panic("error to encode secret")
	}
}

func (controller *SecretController) GetSecret(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	secret := controller.secretService.GetSecret(id)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(secret)
	if err != nil {
		panic("error to encode secret")
	}
}
