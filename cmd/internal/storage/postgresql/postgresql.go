package postgresql

import (
	"avito_go_api/cmd/internal/storage"
	"database/sql"
	"fmt"
	"github.com/lib/pq"

	_ "github.com/lib/pq" // init postgresql driver
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

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
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    segment_id INTEGER REFERENCES segments(id) ON DELETE CASCADE,
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

func (s *Storage) CreateSegment(name string) error {
	const op = "storage.postgres.CreateSegment"
	stmt, err := s.db.Prepare("INSERT INTO segments(name) VALUES($1)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(name)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteSegment(name string) error {
	const op = "storage.postgres.DeleteSegment"
	stmt, err := s.db.Prepare("DELETE FROM segments WHERE name = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) ReassignSegments(addSegments []string, removeSegments []string, userId int64) error {
	const op = "storage.postgres.ReassignSegments"

	stmt, err := s.db.Prepare("DELETE FROM user_segments WHERE user_id = $1 AND segment_id = ANY(SELECT id FROM segments WHERE name = ANY($2))")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(userId, pq.Array(removeSegments))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err = s.db.Prepare("INSERT INTO user_segments(user_id, segment_id) SELECT $1, id FROM segments WHERE name = ANY($2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(userId, pq.Array(addSegments))
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserHasSegment)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetSegments(userId int64) ([]string, error) {
	const op = "storage.postgres.GetSegments"
	stmt, err := s.db.Prepare("SELECT name FROM segments JOIN user_segments ON segments.id = user_segments.segment_id AND user_id = $1")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var segments []string

	for rows.Next() {
		var segmentName string
		if err := rows.Scan(&segmentName); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		segments = append(segments, segmentName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return segments, nil
}
