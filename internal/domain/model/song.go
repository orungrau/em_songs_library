package model

import "time"

type Song struct {
	ID          *string
	Title       *string
	Text        *string
	Link        *string
	Group       *string
	ReleaseDate *time.Time

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

type SongFilter struct {
	ReleaseDateFrom *time.Time
	ReleaseDateTo   *time.Time
	Title           *string
	Text            *string
	Link            *string
	Group           *string
	Page            int
	PageSize        int
}
