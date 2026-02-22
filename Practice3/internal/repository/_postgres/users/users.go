package users

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_postgres "golang/internal/repository/_postgres"
	"golang/pkg/modules"
)

type Repository struct {
	db               *_postgres.Dialect
	executionTimeout time.Duration
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{
		db:               db,
		executionTimeout: time.Second * 5,
	}
}

func (r *Repository) GetUsers(limit, offset int) ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users,
		"SELECT id, name, email, age, role, created_at FROM users WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user,
		"SELECT id, name, email, age, role, password, deleted_at, created_at FROM users WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateUser(req modules.CreateUserRequest) (int, error) {
	var id int
	err := r.db.DB.QueryRow(
		"INSERT INTO users (name, email, age, role, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		req.Name, req.Email, req.Age, req.Role, req.Password,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (r *Repository) UpdateUser(id int, req modules.UpdateUserRequest) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	argIdx := 1
	setClauses := []string{}

	if req.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIdx))
		args = append(args, *req.Name)
		argIdx++
	}
	if req.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, *req.Email)
		argIdx++
	}
	if req.Age != nil {
		setClauses = append(setClauses, fmt.Sprintf("age = $%d", argIdx))
		args = append(args, *req.Age)
		argIdx++
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", argIdx)
	args = append(args, id)

	result, err := r.db.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found or already deleted", id)
	}

	return nil
}

func (r *Repository) DeleteUserByID(id int) (int64, error) {
	result, err := r.db.DB.Exec(
		"UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return 0, fmt.Errorf("user with id %d not found or already deleted", id)
	}

	return rowsAffected, nil
}

func (r *Repository) CreateUserWithAudit(req modules.CreateUserRequest, auditAction, auditDetails string) (int, error) {
	tx, err := r.db.DB.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow(
		"INSERT INTO users (name, email, age, role, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		req.Name, req.Email, req.Age, req.Role, req.Password,
	).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user in transaction: %w", err)
	}

	_, err = tx.Exec(
		"INSERT INTO audit_logs (user_id, action, details) VALUES ($1, $2, $3)",
		userID, auditAction, auditDetails,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create audit log in transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

func (r *Repository) GetUserByEmail(email string) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user,
		"SELECT id, name, email, age, role, password, deleted_at, created_at FROM users WHERE email = $1 AND deleted_at IS NULL", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CountUsers() (int, error) {
	var count int
	err := r.db.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return 0, err
	}
	return count, nil
}
