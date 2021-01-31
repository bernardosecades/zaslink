package http

import (
	"github.com/bernardosecades/sharesecret/service"
	"github.com/bernardosecades/sharesecret/types"

	"github.com/gorilla/mux"

	"encoding/json"
	"net/http"
	"os"
)

type SecretHandler struct {
	secretService service.SecretService
}

func NewSecretHandler(s service.SecretService) *SecretHandler {
	return &SecretHandler{s}
}

func (controller *SecretHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var cs types.CreateSecretRequest
	_ = json.NewDecoder(r.Body).Decode(&cs)

	if len(cs.Content) == 0 {
		customError := types.ErrorResponse{StatusCode: 400, Err: "Empty 'content'"}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	pass := r.Header.Get("X-Password")
	secret, err := controller.secretService.CreateSecret(cs.Content, pass)

	if err != nil {
		customError := types.ErrorResponse{StatusCode: 404, Err: err.Error()}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	cr := types.CreateSecretResponse{
		URL: os.Getenv("SERVER_URL") + ":" + os.Getenv("SERVER_PORT") + "/secret/" + secret.ID,
	}
	_ = json.NewEncoder(w).Encode(cr)
}

func (controller *SecretHandler) GetSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	pass := r.Header.Get("X-Password")

	customPwd, err := controller.secretService.HasSecretWithCustomPwd(id)

	if err != nil {
		customError := types.ErrorResponse{StatusCode: 404, Err: err.Error()}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	if customPwd && pass == "" {
		customError := types.ErrorResponse{StatusCode: 400, Err: "Missing X-Password"}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	content, _ := controller.secretService.GetContentSecret(id, pass)

	if len(content) == 0 {
		customError := types.ErrorResponse{StatusCode: 400, Err: "We can not decrypt the content"}
		w.WriteHeader(customError.StatusCode)
		_ = json.NewEncoder(w).Encode(customError)
		return
	}

	sr := types.SecretResponse{content}
	_ = json.NewEncoder(w).Encode(sr)
}
