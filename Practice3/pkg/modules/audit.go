package modules

import "time"

type AuditLog struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Action    string    `db:"action" json:"action"`
	Details   string    `db:"details" json:"details"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
