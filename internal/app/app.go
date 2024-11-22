package app

import (
	"github.com/orungrau/em_song_library/internal/config"
	"github.com/orungrau/em_song_library/internal/domain/service"
	"github.com/orungrau/em_song_library/internal/repository/storage/song"
	"github.com/orungrau/em_song_library/internal/transport/http"
	"github.com/orungrau/em_song_library/internal/transport/http/handlers"
	"github.com/orungrau/em_song_library/pkg/logger"
	"github.com/orungrau/em_song_library/pkg/transport"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

func MustRun(cfg *config.AppConfig) {
	// Config Logger
	// TODO: Use config
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().
		Timestamp().
		Str("role", "song_library_service").
		Logger()
	log = log.Hook(logger.TracingHook{})
	log.Debug().Msg("Service startup")

	// Config song storage
	songStorage := song.NewPostgresStorage(log, &cfg.PostgresConfig)
	songStorage.MustConnect()
	songStorage.MustAutoMigrate()

	// Config song service
	songService := service.NewSongService(log, songStorage)

	// Config song handler
	songHandler := handlers.NewSongHandler(songService)

	// Setup router
	router := http.NewRouter(log, songHandler, cfg.HttpServer.GetAddress())

	// Start server
	server := transport.NewHTTPServer(log, router, &cfg.HttpServer)
	server.MustStart()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	server.Stop()
	songStorage.Close()
}
