package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // init postgresql driver
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Users table
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY, 
    name varchar(50) NOT NULL
  )`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Users index
	stmt, err = db.Prepare(`CREATE INDEX IF NOT EXISTS idx_users_id ON users(id)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Segments table
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS segments (
    id serial PRIMARY KEY,
    name varchar(50) UNIQUE NOT NULL
  )`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Segments index
	stmt, err = db.Prepare(`CREATE INDEX IF NOT EXISTS idx_segments_name ON segments(name)`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// User segments table
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS user_segments (
    user_id INTEGER REFERENCES users(id),
    segment_id INTEGER REFERENCES segments(id),
    PRIMARY KEY (user_id, segment_id)
  )`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// User segments indexes

	stmt, err = db.Prepare(`CREATE INDEX IF NOT EXISTS idx_user_segments_user_id ON user_segments(user_id)`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = db.Prepare(`CREATE INDEX IF NOT EXISTS idx_user_segments_segment_id ON user_segments(segment_id)`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil

}
