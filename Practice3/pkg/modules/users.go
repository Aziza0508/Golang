package modules

import "time"

type User struct {
	ID        int        `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Email     string     `db:"email" json:"email"`
	Age       int        `db:"age" json:"age"`
	Role      string     `db:"role" json:"role"`
	Password  string     `db:"password" json:"-"`
	DeletedAt *time.Time `db:"deleted_at" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Age   *int    `json:"age,omitempty"`
}
