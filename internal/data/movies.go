package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"forum/internal/validator"
)

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    string    `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
type MovieModel struct {
	DB *sql.DB
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	// v.Check(movie.Genres != nil, "genres", "must be provided")
	// v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	// v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// // Note that we're using the Unique helper in the line below to check that all
	// // values in the movie.Genres slice are unique.
	// v.Check(Unique(movie.Genres), "genres", "must not contain duplicate values")
	// Use the Valid() method to see if any of the checks failed. If they did, then use
	// the failedValidationResponse() helper to send a response to the client, passing
	// in the v.Errors map.
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
func (m MovieModel) Insert(movie *Movie) error {
	// Define the SQL query for inserting a new record in the movies table and returning
	// the system generated data
	stmt := `INSERT INTO movies (title, year,runtime,genres)
	VALUES(?,?,?,?)
	RETURNING id,version`
	// Create an args slice containing the values for the placeholder parameters from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []any{movie.Title, movie.Year, movie.Runtime, movie.Genres}

	// _, err := m.DB.Exec(stmt, movie.Title, movie.Year, movie.Runtime, movie.Genres)
	// if err != nil {
	// 	return err
	// }
	// return nil
	return m.DB.QueryRow(stmt, args...).Scan(&movie.ID, &movie.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Define the sql query for retrieving the movie data
	query := `
	Select id,created_at,title,year,runtime,genres,version
	FROM movies
	WHERE id = ?`
	// Declare a Movie struct to hold the data returned by the query
	var movie Movie
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		&movie.Genres,
		&movie.Version,
	)
	// Handle any errors. If there was no matching movie found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// OtherWise, return a pointer to the Movie struct
	return &movie, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	// Add the 'AND version = $6' clause to the SQL query.
	query := `
	UPDATE movies
	SET title = ?, year = ?, runtime = ?, genres = ?, version = version + 1
	WHERE id = ? AND version = ?
	RETURNING version`
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		movie.Genres,
		movie.ID,
		movie.Version, // Add the expected movie version.
	}
	// Execute the SQL query. If no matching row could be found, we know the movie
	// version has changed (or the record has been deleted) and we return our custom
	// ErrEditConflict error.
	err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
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

// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct te SQL query to delete the record
	query := `
	DELETE FROM movies 
	WHERE id = ?`
	// Execute the SQL query using the Exec() mehtod, passing i9n the id variable as
	// tge value for the placeholder parametr. The ExeC() method returns a sql.Result
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m MovieModel) GetAll(title, genres string, filter Filters) ([]*Movie, error) {
	// Update the sql query to include the filter conditions
	query := `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE (title LIKE '%' || ?1 || '%' OR ?1 = '')
	  AND (genres LIKE '%' || ?2 || '%' OR ?2 = '')
	ORDER BY id ASC `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Pass  the title and genres as the placeholder parametr values
	rows, err := m.DB.QueryContext(ctx, query, title, genres)
	if err != nil {
		return nil, err
	}
	// WHERE (LOWER(title)=LOWER(?) OR ? ='') AND genres = ? OR ?=''
	defer rows.Close()
	movies := []*Movie{}
	for rows.Next() {
		var movie Movie
		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			&movie.Genres,
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}
