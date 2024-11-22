package handlers

import (
	"github.com/orungrau/em_song_library/internal/transport/http/utils"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	utils.Write(w, []byte("OK"))
}
