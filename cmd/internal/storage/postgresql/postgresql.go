package postgresql

import (
	"avito_go_api/cmd/internal/storage"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"os"
	"time"

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

	// User history table
	stmt, err = db.Prepare(`CREATE TABLE IF NOT EXISTS history (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    segment_id INTEGER REFERENCES segments(id) ON DELETE CASCADE,
    operation VARCHAR(10),
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);`)

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

func (s *Storage) RemoveSegmentsFromUser(removeSegments []string, userId int64) error {
	const op = "storage.postgres.RemoveSegmentsFromUser"
	// Removing records from user_segments table
	stmt, err := s.db.Prepare("DELETE FROM user_segments WHERE user_id = $1 AND segment_id = ANY(SELECT id FROM segments WHERE name = ANY($2))")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(userId, pq.Array(removeSegments))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Adding records to history table
	stmt, err = s.db.Prepare("INSERT INTO history(user_id, segment_id,operation) SELECT $1, id, 'Removed' FROM segments WHERE name = ANY($2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(userId, pq.Array(removeSegments))
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) AddSegmentsToUser(addSegments []string, userId int64) error {
	const op = "storage.postgres.RemoveSegmentsFromUser"
	// Adding records to user_segments table
	stmt, err := s.db.Prepare("INSERT INTO user_segments(user_id, segment_id) SELECT $1, id FROM segments WHERE name = ANY($2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(userId, pq.Array(addSegments))
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserHasSegment)
		}
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code == "23503" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	// Adding records to history table
	stmt, err = s.db.Prepare("INSERT INTO history(user_id, segment_id,operation) SELECT $1, id, 'Added' FROM segments WHERE name = ANY($2)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(userId, pq.Array(addSegments))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) ReassignSegments(addSegments []string, removeSegments []string, userId int64) error {
	const op = "storage.postgres.ReassignSegments"

	err := s.RemoveSegmentsFromUser(removeSegments, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	err = s.AddSegmentsToUser(addSegments, userId)
	if err != nil {
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

func (s *Storage) GetUserHistory(userId int64, year int, month time.Month) error {
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	rows, err := s.db.Query("SELECT user_id, segments.name, operation, timestamp FROM history JOIN segments ON segment_id = segments.id WHERE timestamp >= $1 AND timestamp <= $2", startOfMonth, endOfMonth)
	if err != nil {
		return err
	}
	defer rows.Close()

	filename := "latest_segment_history_report.csv"

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "идентификатор пользователя;сегмент;операция;дата и время")

	for rows.Next() {
		var userID int
		var segmentName string
		var operation string
		var timestamp time.Time
		if err := rows.Scan(&userID, &segmentName, &operation, &timestamp); err != nil {
			return err
		}
		fmt.Fprintf(file, "%d;%s;%s;%s\n", userID, segmentName, operation, timestamp.Format("2006-01-02 15:04:05"))
	}

	return nil
}
