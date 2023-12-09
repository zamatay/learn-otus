package domain

import "time"

type Event struct {
	ID           int64         `db:"id"`
	Title        string        `db:"title"`
	Date         time.Time     `db:"date"`
	DateInterval time.Duration `db:"date_interval"`
	Description  string        `db:"description"`
	UserID       int           `db:"user_id"`
	//LeftTime     time.Duration
}
