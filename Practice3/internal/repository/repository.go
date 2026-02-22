package repository

import (
	_postgres "golang/internal/repository/_postgres"
	"golang/internal/repository/_postgres/users"
	"golang/pkg/modules"
)

type UserRepository interface {
	GetUsers(limit, offset int) ([]modules.User, error)
	GetUserByID(id int) (*modules.User, error)
	CreateUser(req modules.CreateUserRequest) (int, error)
	UpdateUser(id int, req modules.UpdateUserRequest) error
	DeleteUserByID(id int) (int64, error)
	CreateUserWithAudit(req modules.CreateUserRequest, auditAction, auditDetails string) (int, error)
	GetUserByEmail(email string) (*modules.User, error)
	CountUsers() (int, error)
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
