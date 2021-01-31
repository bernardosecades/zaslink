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

func (controller *SecretHandler) CreateSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var cs types.CreateSecretRequest
	err := json.NewDecoder(r.Body).Decode(&cs)

	if err != nil {
		customError := types.ErrorResponse{StatusCode: 400, Err: "Bad JSON in CreateSecretRequest'"}
		encodeCustomError(customError, w)
		return
	}

	if len(cs.Content) == 0 {
		customError := types.ErrorResponse{StatusCode: 400, Err: "Empty 'content'"}
		encodeCustomError(customError, w)
		return
	}

	pass := r.Header.Get("X-Password")
	secret, err := controller.secretService.CreateSecret(cs.Content, pass)

	if err != nil {
		customError := types.ErrorResponse{StatusCode: 500, Err: "We can not create your secret"}
		log.Println(err.Error())
		encodeCustomError(customError, w)
		return
	}

	cr := types.CreateSecretResponse{
		URL: os.Getenv("SERVER_URL") + ":" + os.Getenv("SERVER_PORT") + "/secret/" + secret.ID,
	}

	if err := json.NewEncoder(w).Encode(cr); err != nil {
		log.Fatal("Error to encode CreateSecretResponse")
	}
}

func (controller *SecretHandler) GetSecret(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	pass := r.Header.Get("X-Password")

	customPwd, err := controller.secretService.HasSecretWithCustomPwd(id)

	if err != nil {
		customError := types.ErrorResponse{StatusCode: 404, Err: "Unknown Secret: It either never existed or has already been viewed"}
		log.Println(err.Error())
		encodeCustomError(customError, w)
		return
	}

	if customPwd && pass == "" {
		customError := types.ErrorResponse{StatusCode: 400, Err: "Missing X-Password"}
		encodeCustomError(customError, w)
		return
	}

	content, err := controller.secretService.GetContentSecret(id, pass)

	if err != nil {
		customError := types.ErrorResponse{StatusCode: 404, Err: "Unknown Secret: It either never existed or has already been viewed"}
		log.Println(err.Error())
		encodeCustomError(customError, w)
		return
	}

	if len(content) == 0 {
		customError := types.ErrorResponse{StatusCode: 400, Err: "We can not decrypt the content"}
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
