package repository

import (
	"assignment4/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]models.Movie, error) {
	rows, err := r.db.Query(ctx, "SELECT id, genre, budget, title, hero, heroine, created_at FROM movies ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("query movies: %w", err)
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Genre, &m.Budget, &m.Title, &m.Hero, &m.Heroine, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan movie: %w", err)
		}
		movies = append(movies, m)
	}
	if movies == nil {
		movies = []models.Movie{}
	}
	return movies, rows.Err()
}

func (r *MovieRepository) GetByID(ctx context.Context, id int) (models.Movie, error) {
	var m models.Movie
	err := r.db.QueryRow(ctx,
		"SELECT id, genre, budget, title, hero, heroine, created_at FROM movies WHERE id = $1", id,
	).Scan(&m.ID, &m.Genre, &m.Budget, &m.Title, &m.Hero, &m.Heroine, &m.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("get movie %d: %w", id, err)
	}
	return m, nil
}

func (r *MovieRepository) Create(ctx context.Context, req models.CreateMovieRequest) (models.Movie, error) {
	var m models.Movie
	err := r.db.QueryRow(ctx,
		`INSERT INTO movies (genre, budget, title, hero, heroine)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, genre, budget, title, hero, heroine, created_at`,
		req.Genre, req.Budget, req.Title, req.Hero, req.Heroine,
	).Scan(&m.ID, &m.Genre, &m.Budget, &m.Title, &m.Hero, &m.Heroine, &m.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("create movie: %w", err)
	}
	return m, nil
}

func (r *MovieRepository) Update(ctx context.Context, id int, req models.UpdateMovieRequest) (models.Movie, error) {
	var m models.Movie
	err := r.db.QueryRow(ctx,
		`UPDATE movies SET
			genre    = COALESCE($1, genre),
			budget   = COALESCE($2, budget),
			title    = COALESCE($3, title),
			hero     = COALESCE($4, hero),
			heroine  = COALESCE($5, heroine)
		 WHERE id = $6
		 RETURNING id, genre, budget, title, hero, heroine, created_at`,
		req.Genre, req.Budget, req.Title, req.Hero, req.Heroine, id,
	).Scan(&m.ID, &m.Genre, &m.Budget, &m.Title, &m.Hero, &m.Heroine, &m.CreatedAt)
	if err != nil {
		return models.Movie{}, fmt.Errorf("update movie %d: %w", id, err)
	}
	return m, nil
}

func (r *MovieRepository) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete movie %d: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("movie %d not found", id)
	}
	return nil
}
