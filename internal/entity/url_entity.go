package entity

import "time"

type URL struct {
	ID         uint64    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`
	URL        string    `db:"url" json:"url"`
	ShortCode  string    `db:"short_code" json:"short_code"`
	VisitCount uint64    `db:"visit_count" json:"visit_count"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
