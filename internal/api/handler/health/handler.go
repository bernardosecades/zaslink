package health

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// TODO move trace and metric
var (
	tracer      = otel.Tracer("sharesecret-service")
	meter       = otel.Meter("sharesecret-service")
	viewCounter metric.Int64Counter
)

func init() {
	var err error
	viewCounter, err = meter.Int64Counter("user.views",
		metric.WithDescription("The number of views"),
		metric.WithUnit("{views}"))
	if err != nil {
		panic(err)
	}
}

type Handler struct {
	mongodbURI string
}

func NewHandler(mongodbURI string) *Handler {
	return &Handler{mongodbURI: mongodbURI}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	// Note: if the context already contains an span then newly Span will be child of that span
	ctx, span := tracer.Start(r.Context(), "Healthz")
	defer span.End()

	viewCounter.Add(ctx, 1)

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
