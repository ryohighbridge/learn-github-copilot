package repository

import (
	"database/sql"
	"time"

	"github.com/ryohighbridge/learn-github-copilot/backend/internal/domain"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// GetAll 全てのイベントを取得
func (r *EventRepository) GetAll() ([]domain.Event, error) {
	query := `SELECT id, title, description, start_date, end_date, all_day, created_at, updated_at 
	          FROM events ORDER BY start_date ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.AllDay,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// GetByID IDでイベントを取得
func (r *EventRepository) GetByID(id int) (*domain.Event, error) {
	query := `SELECT id, title, description, start_date, end_date, all_day, created_at, updated_at 
	          FROM events WHERE id = $1`

	var event domain.Event
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartDate,
		&event.EndDate,
		&event.AllDay,
		&event.CreatedAt,
		&event.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// GetByDateRange 期間内のイベントを取得
func (r *EventRepository) GetByDateRange(start, end time.Time) ([]domain.Event, error) {
	query := `SELECT id, title, description, start_date, end_date, all_day, created_at, updated_at 
	          FROM events 
	          WHERE start_date <= $2 AND end_date >= $1
	          ORDER BY start_date ASC`

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.StartDate,
			&event.EndDate,
			&event.AllDay,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// Create 新しいイベントを作成
func (r *EventRepository) Create(event *domain.Event) error {
	query := `INSERT INTO events (title, description, start_date, end_date, all_day) 
	          VALUES ($1, $2, $3, $4, $5) 
	          RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.AllDay,
	).Scan(&event.ID, &event.CreatedAt, &event.UpdatedAt)
}

// Update イベントを更新
func (r *EventRepository) Update(event *domain.Event) error {
	query := `UPDATE events 
	          SET title = $1, description = $2, start_date = $3, end_date = $4, all_day = $5
	          WHERE id = $6
	          RETURNING updated_at`

	return r.db.QueryRow(
		query,
		event.Title,
		event.Description,
		event.StartDate,
		event.EndDate,
		event.AllDay,
		event.ID,
	).Scan(&event.UpdatedAt)
}

// Delete イベントを削除
func (r *EventRepository) Delete(id int) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
