package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/orungrau/em_song_library/docs"
	"github.com/orungrau/em_song_library/internal/transport/http/handlers"
	"github.com/orungrau/em_song_library/internal/transport/http/middleware"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func NewRouter(log zerolog.Logger, songHandler *handlers.SongHandler, address string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.NewLoggerMiddleware(log).Middleware)

	r.Get("/health", handlers.HealthCheck)

	r.Route("/songs", func(r chi.Router) {
		r.Get("/", songHandler.GetAll)
		r.Get("/{songId}", songHandler.Get)
		r.Post("/", songHandler.Create)
		r.Patch("/{songId}", songHandler.Update)
		r.Delete("/{songId}", songHandler.Delete)
	})

	log.Debug().Msg(fmt.Sprintf("Swagger available at http://%s/swagger/index.html", address))
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", address)), // The url pointing to API definition
	))

	return r
}
