package grpc

import (
	"context"
	sharesecretgrpc "github.com/bernardosecades/sharesecret/build"
	"github.com/bernardosecades/sharesecret/service"
	_ "google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/status"
)


type shareSecretHandler struct {
	secretService service.SecretService
}

func NewShareSecretServer(s service.SecretService) sharesecretgrpc.SecretAppServer {
	return &shareSecretHandler{secretService: s}
}

func (s shareSecretHandler) CreateSecret(ctx context.Context, req *sharesecretgrpc.CreateSecretRequest) (*sharesecretgrpc.CreateSecretResponse, error) {

	secret, err := s.secretService.CreateSecret(req.Content, req.Password)

	// TODO handle errors
	/*
	return data, status.Errorf(
	            codes.InvalidArgument,
	            fmt.Sprintf("Your message", req.data),
	        )
	 */

	if err != nil {
		return nil, err
	}

	r := &sharesecretgrpc.CreateSecretResponse{}
	r.Id = secret.ID

	return r, nil
}

func (s shareSecretHandler) SeeSecret(ctx context.Context, req *sharesecretgrpc.SeeSecretRequest) (*sharesecretgrpc.SeeSecretResponse, error) {
	c, err :=s.secretService.GetContentSecret(req.Id, req.Password)

	if err != nil {
		return nil, err
	}

	r := &sharesecretgrpc.SeeSecretResponse{}
	r.Content = c

	return r, nil
}
