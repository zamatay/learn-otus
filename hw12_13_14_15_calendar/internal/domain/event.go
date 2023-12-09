package domain

import "time"

type Event struct {
	Id           int64         `db:"id"`
	Title        string        `db:"title"`
	Date         time.Time     `db:"date"`
	DateInterval time.Duration `db:"date_interval"`
	Description  string        `db:"description"`
	UserId       int           `db:"user_id"`
	//LeftTime     time.Duration
}
