package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"greenlight/internal/validator"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

func (m Movie) MarshalJSON() ([]byte, error) {
	var runtime string

	if m.Runtime != 0 {
		runtime = fmt.Sprintf("%d mins", m.Runtime)
	}

	aux := struct {
		ID      int64    `json:"id"`
		Title   string   `json:"title"`
		Year    int32    `json:"year,omitempty"`
		Runtime string   `json:"runtime,omitempty"` // This is a string.
		Genres  []string `json:"genres,omitempty"`
		Version int32    `json:"version"`
	}{
		// Set the values for the anonymous struct.
		ID:      m.ID,
		Title:   m.Title,
		Year:    m.Year,
		Runtime: runtime, // Note that we assign the value from the runtime variable here.
		Genres:  m.Genres,
		Version: m.Version,
	}

	return json.Marshal(aux)
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

// Insert method for inserting a new record in the movies table.
func (m MovieModel) Insert(movie *Movie) error {
	query := `
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Remove the pg_sleep(10) clause.
	query := `
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1`

	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Remove &[]byte{} from the first Scan() destination.
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	query := `
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM movies
WHERE id = $1`

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Create a new GetAll() method which returns a slice of movies. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	query := fmt.Sprintf(`
SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY %s %s, id ASC
LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []interface{}{title, pq.Array(genres), filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return movies, metadata, nil
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {
	// Mock the action...
	return nil
}

func (m MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MockMovieModel) Update(movie *Movie) error {
	// Mock the action...
	return nil
}

func (m MockMovieModel) Delete(id int64) error {
	// Mock the action...
	return nil
}
