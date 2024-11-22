package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/orungrau/em_song_library/internal/domain/model"
	"github.com/orungrau/em_song_library/internal/domain/service"
	"github.com/orungrau/em_song_library/internal/transport/http/dto"
	"github.com/orungrau/em_song_library/internal/transport/http/utils"
	"net/http"
)

type SongHandler struct {
	validate    *validator.Validate
	decoder     *schema.Decoder
	songService service.SongService
}

func NewSongHandler(songService service.SongService) *SongHandler {
	validate := validator.New()
	decoder := schema.NewDecoder()

	decoder.IgnoreUnknownKeys(true)
	decoder.ZeroEmpty(true)

	return &SongHandler{
		decoder:     decoder,
		validate:    validate,
		songService: songService,
	}
}

// GetAll godoc
// @Summary Get all songs
// @Description Retrieve a list of songs with optional filters such as release date range, title, group, and pagination.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param release_date_from query int64 false "Filter songs by release date from (Unix timestamp)"
// @Param release_date_to query int64 false "Filter songs by release date to (Unix timestamp)"
// @Param title query string false "Filter songs by title"
// @Param text query string false "Filter songs by text"
// @Param link query string false "Filter songs by link"
// @Param group query string false "Filter songs by group"
// @Param page query int false "Page number (default: 0)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {object} dto.SongList "A paginated list of songs"
// @Failure 400 {object} dto.Status "Bad request error with a detailed message"
// @Router /songs [get]
func (h *SongHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.WriteErrorJson(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	var filter dto.SongFilter
	err = h.decoder.Decode(&filter, r.Form)
	if err != nil {
		utils.WriteErrorJson(w, "Failed to decode filter: "+err.Error(), http.StatusBadRequest)
		return
	}

	var songs []dto.Song = make([]dto.Song, 0)
	total := len(songs)

	response := dto.SongList{
		Data:       songs,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: (total + filter.PageSize - 1) / filter.PageSize, // Calculate total pages
	}

	utils.WriteJson(w, response, http.StatusOK)
}

// Get godoc
// @Summary Get a single song
// @Description Retrieve details of a specific song by its ID.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param songId path string true "ID of the song to retrieve"
// @Success 200 {object} dto.Song "Details of the requested song"
// @Failure 400 {object} dto.Status "Bad request error with a detailed message"
// @Router /songs/{songId} [get]
func (h *SongHandler) Get(w http.ResponseWriter, r *http.Request) {
	songId := chi.URLParam(r, "songId")

	song, err := h.songService.Get(r.Context(), songId)
	if err != nil {
		utils.WriteErrorJson(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if song == nil {
		utils.WriteErrorJson(w, "Not found", http.StatusNotFound)
		return
	}

	createdSong := dto.Song{
		ID:          *song.ID,
		Title:       *song.Title,
		Group:       *song.Group,
		ReleaseDate: *song.ReleaseDate,
	}

	utils.WriteJson(w, createdSong, http.StatusOK)
}

// Create godoc
// @Summary Create a new song
// @Description Add a new song to the library by providing required details.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body dto.CreateSong true "Details of the song to create"
// @Success 201 {object} dto.Song "The created song"
// @Failure 400 {object} dto.Status "Bad request error with a detailed message"
// @Router /songs [post]
func (h *SongHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createDTO dto.CreateSong

	if err := json.NewDecoder(r.Body).Decode(&createDTO); err != nil {
		utils.WriteErrorJson(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(createDTO); err != nil {
		utils.WriteErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdSong, err := h.songService.Create(r.Context(), model.Song{
		Title:       &createDTO.Title,
		Group:       &createDTO.Group,
		ReleaseDate: &createDTO.ReleaseDate.Time,
	})
	if err != nil {
		utils.WriteErrorJson(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, dto.SongFromModel(createdSong), http.StatusCreated)
}

// Update godoc
// @Summary Update an existing song
// @Description Update the details of an existing song. The song ID must be specified in the request.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param songId path string true "ID of the song to update"
// @Param song body dto.Song true "Updated details of the song"
// @Success 200 {object} dto.Song "The updated song"
// @Failure 400 {object} dto.Status "Bad request error with a detailed message"
// @Router /songs/{songId} [patch]
func (h *SongHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updateDTO dto.Song

	if err := json.NewDecoder(r.Body).Decode(&updateDTO); err != nil {
		utils.WriteErrorJson(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(updateDTO); err != nil {
		utils.WriteErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedSong, err := h.songService.Update(r.Context(), model.Song{
		ID:          &updateDTO.ID,
		Title:       &updateDTO.Title,
		Text:        updateDTO.Text,
		Link:        updateDTO.Link,
		Group:       &updateDTO.Group,
		ReleaseDate: &updateDTO.ReleaseDate,
	})

	if err != nil {
		utils.WriteErrorJson(w, err.Error(), http.StatusBadRequest)
	}

	utils.WriteJson(w, dto.SongFromModel(updatedSong), http.StatusOK)
}

// Delete godoc
// @Summary Delete a song
// @Description Remove a song from the library by its ID.
// @Tags songs
// @Accept  json
// @Produce  json
// @Param songId path string true "ID of the song to delete"
// @Success 200 {object} dto.Status "Confirmation of successful deletion"
// @Failure 400 {object} dto.Status "Bad request error with a detailed message"
// @Router /songs/{songId} [delete]
func (h *SongHandler) Delete(w http.ResponseWriter, r *http.Request) {
	songId := chi.URLParam(r, "songId")

	err := h.songService.Delete(r.Context(), songId)
	if err != nil {
		utils.WriteErrorJson(w, err.Error(), http.StatusBadRequest)
	}

	status := dto.Status{
		Error:   false,
		Message: fmt.Sprintf("song deleted with id: %s", songId),
	}

	utils.WriteJson(w, status, http.StatusOK)
}
