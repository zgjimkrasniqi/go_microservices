package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) LogEntry {
	db = dbPool
	return LogEntry{}
}

type LogEntry struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	Data      string `json:"data"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (l *LogEntry) Insert(logEntry LogEntry) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into logger (name, data, created_at, updated_at)
		values ($1, $2, $3, $4) returning id`

	err := db.QueryRowContext(ctx, stmt,
		logEntry.Name,
		logEntry.Data,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (l *LogEntry) GetAllLogs() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, name, data, created_at, updated_at
	from logger order by id`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*LogEntry

	for rows.Next() {
		var logEntry LogEntry
		err := rows.Scan(
			&logEntry.ID,
			&logEntry.Name,
			&logEntry.Data,
			&logEntry.CreatedAt,
			&logEntry.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		logs = append(logs, &logEntry)
	}

	return logs, nil
}
