package main

import (
	"github.com/orungrau/em_song_library/internal/app"
	"github.com/orungrau/em_song_library/internal/config"
)

func main() {
	cfg := config.MustLoad()
	app.MustRun(cfg)
}
