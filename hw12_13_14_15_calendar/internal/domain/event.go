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

func NewEvent(ID int64, UserID int32, Title string, Description string, DateInterval int32, Date uint64) *Event {
	return &Event{ID: ID, UserID: int(UserID), Title: Title, Description: Description, Date: time.Unix(int64(Date), 0), DateInterval: time.Duration(DateInterval)}
}
