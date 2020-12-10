package controllers

import (
	"encoding/json"
	"github.com/bernardosecades/sharesecret/services"
	models "github.com/bernardosecades/sharesecret/viewmodel"
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
	secret := controller.secretService.CreateSecret("raw content", "myPassword")

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(secret)
	if err != nil {
		panic("error to encode secret")
	}
}

func (controller *SecretController) GetSecret(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pass:= r.Header.Get("X-Password")

	if pass == "" {
		pass = os.Getenv("SECRET_PASSWORD")
	}

	content := controller.secretService.GetContentSecret(id, pass)
	viewModel := models.SecretViewModel{content}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(viewModel)
	if err != nil {
		panic("error to encode secret")
	}
}
