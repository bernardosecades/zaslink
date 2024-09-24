package health

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Handler struct {
	mongodbURI string
}

func NewHandler(mongodbURI string) *Handler {
	return &Handler{mongodbURI: mongodbURI}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	opts := options.Client().ApplyURI(h.mongodbURI)
	client, _ := mongo.Connect(opts)
	err := client.Ping(r.Context(), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
