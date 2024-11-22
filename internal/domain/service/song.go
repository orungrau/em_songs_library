package service

import (
	"context"
	"github.com/orungrau/em_song_library/internal/domain/model"
	"github.com/rs/zerolog"
)

type SongStorage interface {
	GetByFilters(ctx context.Context, filters model.SongFilter) ([]*model.Song, error)
	GetById(ctx context.Context, id string, allowDeleted bool) (*model.Song, error)
	Create(ctx context.Context, song model.Song) (*model.Song, error)
	Update(ctx context.Context, song model.Song) (*model.Song, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	DeletePermanent(ctx context.Context, id string) error
}

type SongService interface {
	GetByFilters(ctx context.Context, filters model.SongFilter) ([]*model.Song, error)
	Get(ctx context.Context, id string) (*model.Song, error)
	Create(ctx context.Context, song model.Song) (*model.Song, error)
	Update(ctx context.Context, song model.Song) (*model.Song, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	DeletePermanent(ctx context.Context, id string) error
}

type songService struct {
	log     zerolog.Logger
	storage SongStorage
}

func NewSongService(log zerolog.Logger, storage SongStorage) SongService {
	return &songService{
		log:     log.With().Str("module", "song-service").Logger(),
		storage: storage,
	}
}

func (s *songService) GetByFilters(ctx context.Context, filters model.SongFilter) ([]*model.Song, error) {
	return s.storage.GetByFilters(ctx, filters)
}

func (s *songService) Get(ctx context.Context, id string) (*model.Song, error) {
	return s.storage.GetById(ctx, id, false)
}

func (s *songService) Create(ctx context.Context, song model.Song) (*model.Song, error) {
	return s.storage.Create(ctx, song)
}

func (s *songService) Update(ctx context.Context, song model.Song) (*model.Song, error) {
	return s.storage.Update(ctx, song)
}

func (s *songService) Delete(ctx context.Context, id string) error {
	return s.storage.Delete(ctx, id)
}

func (s *songService) Restore(ctx context.Context, id string) error {
	return s.storage.Restore(ctx, id)
}

func (s *songService) DeletePermanent(ctx context.Context, id string) error {
	return s.storage.DeletePermanent(ctx, id)
}
