package song

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/orungrau/em_song_library/internal/domain/model"
	"github.com/orungrau/em_song_library/internal/domain/service"
	"github.com/rs/zerolog"
)

type PostgresStorageConfig interface {
	GetConnection() string
	GetDatabase() string
	GetMigrationSource() string
}

type PostgresStorage interface {
	service.SongStorage
	MustAutoMigrate()
	MustConnect()
	Close()
}

type songPostgresStorage struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
	cfg  PostgresStorageConfig
}

func NewPostgresStorage(log zerolog.Logger, cfg PostgresStorageConfig) PostgresStorage {
	return &songPostgresStorage{
		cfg: cfg,
		log: log.With().Str("module", "song-postgres-storage").Logger(),
	}
}

func (s *songPostgresStorage) MustConnect() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, s.cfg.GetConnection())
	if err != nil {
		s.log.Fatal().Err(err).Msg("Failed to create PostgreSQL connection pool")
		return
	}

	if err = pool.Ping(ctx); err != nil {
		s.log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL database")
		return
	}

	s.log.Info().Msg("Successfully connected to PostgreSQL database")
	s.pool = pool
}

func (s *songPostgresStorage) Close() {
	if s.pool != nil {
		s.pool.Close()
		s.log.Info().Msg("PostgreSQL connection pool closed")
	}
}

func (s *songPostgresStorage) MustAutoMigrate() {
	db := stdlib.OpenDBFromPool(s.pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		s.log.Fatal().Err(err).Msg("Could not create postgres driver")
	}

	defer func() {
		if err := driver.Close(); err != nil {
			s.log.Fatal().Err(err).Msg("Error closing migrations connection")
		}
	}()

	m, err := migrate.NewWithDatabaseInstance(
		s.cfg.GetMigrationSource(),
		s.cfg.GetDatabase(),
		driver,
	)
	if err != nil {
		s.log.Fatal().Err(err).Msg("Could not create migrate instance")
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			s.log.Info().Msg("No migrations to apply")
		} else {
			s.log.Fatal().Err(err).Msg("Could not apply migrations")
		}
		return
	}

	s.log.Info().Msg("Migrations applied successfully")
}

func (s *songPostgresStorage) GetByFilters(ctx context.Context, filters model.SongFilter) ([]*model.Song, error) {
	query := `
		SELECT id, title, text, link, "group", release_date, created_at, updated_at, deleted_at
		FROM songs
		WHERE 1=1`
	var args []interface{}
	argIndex := 1

	if filters.ReleaseDateFrom != nil {
		query += fmt.Sprintf(" AND release_date >= $%d", argIndex)
		args = append(args, *filters.ReleaseDateFrom)
		argIndex++
	}
	if filters.ReleaseDateTo != nil {
		query += fmt.Sprintf(" AND release_date <= $%d", argIndex)
		args = append(args, *filters.ReleaseDateTo)
		argIndex++
	}
	if filters.Title != nil {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, "%"+*filters.Title+"%")
		argIndex++
	}
	if filters.Text != nil {
		query += fmt.Sprintf(" AND text ILIKE $%d", argIndex)
		args = append(args, "%"+*filters.Text+"%")
		argIndex++
	}
	if filters.Link != nil {
		query += fmt.Sprintf(" AND link ILIKE $%d", argIndex)
		args = append(args, "%"+*filters.Link+"%")
		argIndex++
	}
	if filters.Group != nil {
		query += fmt.Sprintf(" AND \"group\" ILIKE $%d", argIndex)
		args = append(args, "%"+*filters.Group+"%")
		argIndex++
	}

	query += " ORDER BY release_date DESC"

	if filters.PageSize >= 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.PageSize)
		argIndex++
	}
	if filters.Page >= 0 && filters.PageSize > 0 {
		offset := filters.Page * filters.PageSize
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
		argIndex++
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	songs := make([]*model.Song, 0)
	for rows.Next() {
		var song model.Song
		err := rows.Scan(
			&song.ID,
			&song.Title,
			&song.Text,
			&song.Link,
			&song.Group,
			&song.ReleaseDate,
			&song.CreatedAt,
			&song.UpdatedAt,
			&song.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *songPostgresStorage) GetById(ctx context.Context, id string, allowDeleted bool) (*model.Song, error) {
	query := `
		SELECT id, title, text, link, "group", release_date, created_at, updated_at, deleted_at
		FROM songs 
		WHERE id = $1`

	if !allowDeleted {
		query += " AND deleted_at IS NULL"
	}

	var song model.Song
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&song.ID,
		&song.Title,
		&song.Text,
		&song.Link,
		&song.Group,
		&song.ReleaseDate,
		&song.CreatedAt,
		&song.UpdatedAt,
		&song.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &song, nil
}

func (s *songPostgresStorage) Create(ctx context.Context, song model.Song) (*model.Song, error) {
	query := `
		INSERT INTO songs (title, text, link, "group", release_date, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var id string
	err := s.pool.QueryRow(
		ctx,
		query,
		song.Title,
		song.Text,
		song.Link,
		song.Group,
		song.ReleaseDate,
		song.CreatedAt,
		song.UpdatedAt,
		song.DeletedAt,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	song.ID = &id
	return &song, nil
}

func (s *songPostgresStorage) Update(ctx context.Context, song model.Song) (*model.Song, error) {
	if song.ID == nil {
		return nil, fmt.Errorf("ID is required")
	}

	query := `UPDATE songs SET `
	var args []interface{}
	argIndex := 1

	if song.Title != nil {
		query += `title = $` + fmt.Sprint(argIndex) + `, `
		args = append(args, *song.Title)
		argIndex++
	}
	if song.Text != nil {
		query += `"text" = $` + fmt.Sprint(argIndex) + `, `
		args = append(args, *song.Text)
		argIndex++
	}
	if song.Link != nil {
		query += `"link" = $` + fmt.Sprint(argIndex) + `, `
		args = append(args, *song.Link)
		argIndex++
	}
	if song.Group != nil {
		query += `"group" = $` + fmt.Sprint(argIndex) + `, `
		args = append(args, *song.Group)
		argIndex++
	}
	if song.ReleaseDate != nil {
		query += `release_date = $` + fmt.Sprint(argIndex) + `, `
		args = append(args, *song.ReleaseDate)
		argIndex++
	}

	query += `WHERE id = $` + fmt.Sprint(argIndex) + ` AND deleted_at IS NULL 
		RETURNING id, title, text, link, "group", release_date, created_at, updated_at, deleted_at`
	args = append(args, *song.ID)

	var updatedSong model.Song
	err := s.pool.QueryRow(ctx, query, args...).Scan(
		&updatedSong.ID,
		&updatedSong.Title,
		&updatedSong.Text,
		&updatedSong.Link,
		&updatedSong.Group,
		&updatedSong.ReleaseDate,
		&updatedSong.CreatedAt,
		&updatedSong.UpdatedAt,
		&updatedSong.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("song not found or already deleted")
		}
		return nil, err
	}

	return &updatedSong, nil
}

func (s *songPostgresStorage) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE songs 
		SET deleted_at = CURRENT_TIMESTAMP 
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("song not found or already deleted")
	}

	return nil
}

func (s *songPostgresStorage) Restore(ctx context.Context, id string) error {
	query := `
		UPDATE songs 
		SET deleted_at = NULL 
		WHERE id = $1 AND deleted_at IS NOT NULL
	`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("song not found or not deleted")
	}

	return nil
}

func (s *songPostgresStorage) DeletePermanent(ctx context.Context, id string) error {
	query := `
		DELETE FROM songs 
		WHERE id = $1
	`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("song not found")
	}

	return nil
}
