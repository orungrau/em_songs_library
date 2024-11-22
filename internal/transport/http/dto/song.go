package dto

import (
	"github.com/orungrau/em_song_library/internal/domain/model"
	"time"
)

type Song struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Text        *string   `json:"text"`
	Link        *string   `json:"link"`
	Group       string    `json:"group"`
	ReleaseDate time.Time `json:"release_date"`
} // @name Song

func SongFromModel(song *model.Song) Song {
	return Song{
		ID:          *song.ID,
		Title:       *song.Title,
		Group:       *song.Group,
		ReleaseDate: *song.ReleaseDate,
	}
}

func (m *Song) ToModel() model.Song {
	return model.Song{
		ID:          &m.ID,
		Title:       &m.Title,
		Group:       &m.Group,
		ReleaseDate: &m.ReleaseDate,
	}
}

type SongFilter struct {
	ReleaseDateFrom *time.Time `json:"release_date_from" schema:"release_date_from"`
	ReleaseDateTo   *time.Time `json:"release_date_to" schema:"release_date_from"`
	Title           *string    `json:"title"`
	Text            *string    `json:"text"`
	Link            *string    `json:"link"`
	Group           *string    `json:"group"`
	Page            int        `json:"page" schema:"page,default:0"`
	PageSize        int        `json:"page_size" schema:"page_size,default:10"`
}

type CreateSong struct {
	Title       string        `json:"title"`
	Text        *string       `json:"text,omitempty"`
	Link        *string       `json:"link,omitempty" `
	Group       string        `json:"group"`
	ReleaseDate TimestampTime `json:"release_date" swaggertype:"primitive,integer"`
} // @name CreateSong

type SongList struct {
	Data     []Song `json:"data"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
} // @name SongList
