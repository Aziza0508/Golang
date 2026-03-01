package models

import "time"

type Movie struct {
	ID        int       `json:"id"`
	Genre     string    `json:"genre"`
	Budget    int64     `json:"budget"`
	Title     string    `json:"title"`
	Hero      string    `json:"hero"`
	Heroine   string    `json:"heroine"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMovieRequest struct {
	Genre   string `json:"genre"`
	Budget  int64  `json:"budget"`
	Title   string `json:"title"`
	Hero    string `json:"hero"`
	Heroine string `json:"heroine"`
}

type UpdateMovieRequest struct {
	Genre   *string `json:"genre"`
	Budget  *int64  `json:"budget"`
	Title   *string `json:"title"`
	Hero    *string `json:"hero"`
	Heroine *string `json:"heroine"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
