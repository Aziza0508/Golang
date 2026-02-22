package users_test

import (
	"fmt"
	"testing"
	"time"

	"golang/internal/usecase/users"
	"golang/pkg/modules"
)

// --- Mock Repository ---

type MockUserRepository struct {
	GetUsersFunc            func(limit, offset int) ([]modules.User, error)
	GetUserByIDFunc         func(id int) (*modules.User, error)
	CreateUserFunc          func(req modules.CreateUserRequest) (int, error)
	UpdateUserFunc          func(id int, req modules.UpdateUserRequest) error
	DeleteUserByIDFunc      func(id int) (int64, error)
	CreateUserWithAuditFunc func(req modules.CreateUserRequest, action, details string) (int, error)
	GetUserByEmailFunc      func(email string) (*modules.User, error)
	CountUsersFunc          func() (int, error)
}

func (m *MockUserRepository) GetUsers(limit, offset int) ([]modules.User, error) {
	return m.GetUsersFunc(limit, offset)
}

func (m *MockUserRepository) GetUserByID(id int) (*modules.User, error) {
	return m.GetUserByIDFunc(id)
}

func (m *MockUserRepository) CreateUser(req modules.CreateUserRequest) (int, error) {
	return m.CreateUserFunc(req)
}

func (m *MockUserRepository) UpdateUser(id int, req modules.UpdateUserRequest) error {
	return m.UpdateUserFunc(id, req)
}

func (m *MockUserRepository) DeleteUserByID(id int) (int64, error) {
	return m.DeleteUserByIDFunc(id)
}

func (m *MockUserRepository) CreateUserWithAudit(req modules.CreateUserRequest, action, details string) (int, error) {
	return m.CreateUserWithAuditFunc(req, action, details)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*modules.User, error) {
	return m.GetUserByEmailFunc(email)
}

func (m *MockUserRepository) CountUsers() (int, error) {
	return m.CountUsersFunc()
}

// --- Tests ---

func TestGetUsers_Success(t *testing.T) {
	mock := &MockUserRepository{
		GetUsersFunc: func(limit, offset int) ([]modules.User, error) {
			return []modules.User{
				{ID: 1, Name: "John", Email: "john@example.com", Age: 30, Role: "admin", CreatedAt: time.Now()},
				{ID: 2, Name: "Jane", Email: "jane@example.com", Age: 25, Role: "user", CreatedAt: time.Now()},
			}, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	result, err := uc.GetUsers(10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result))
	}
}

func TestGetUsers_EmptyResult(t *testing.T) {
	mock := &MockUserRepository{
		GetUsersFunc: func(limit, offset int) ([]modules.User, error) {
			return []modules.User{}, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	result, err := uc.GetUsers(10, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 users, got %d", len(result))
	}
}

func TestGetUsers_PaginationDefaults(t *testing.T) {
	var capturedLimit, capturedOffset int
	mock := &MockUserRepository{
		GetUsersFunc: func(limit, offset int) ([]modules.User, error) {
			capturedLimit = limit
			capturedOffset = offset
			return []modules.User{}, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)

	// Test negative limit defaults to 10
	_, _ = uc.GetUsers(-1, 0)
	if capturedLimit != 10 {
		t.Fatalf("expected default limit 10, got %d", capturedLimit)
	}

	// Test limit > 100 capped to 100
	_, _ = uc.GetUsers(200, 0)
	if capturedLimit != 100 {
		t.Fatalf("expected max limit 100, got %d", capturedLimit)
	}

	// Test negative offset defaults to 0
	_, _ = uc.GetUsers(10, -5)
	if capturedOffset != 0 {
		t.Fatalf("expected default offset 0, got %d", capturedOffset)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	mock := &MockUserRepository{
		GetUserByIDFunc: func(id int) (*modules.User, error) {
			return &modules.User{ID: 1, Name: "John", Email: "john@example.com", Age: 30, Role: "admin"}, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	user, err := uc.GetUserByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("expected user ID 1, got %d", user.ID)
	}
	if user.Name != "John" {
		t.Fatalf("expected name John, got %s", user.Name)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	mock := &MockUserRepository{
		GetUserByIDFunc: func(id int) (*modules.User, error) {
			return nil, fmt.Errorf("user with id %d not found", id)
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	user, err := uc.GetUserByID(999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if user != nil {
		t.Fatal("expected nil user")
	}
}

func TestCreateUser_Success(t *testing.T) {
	mock := &MockUserRepository{
		CreateUserFunc: func(req modules.CreateUserRequest) (int, error) {
			return 1, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	id, err := uc.CreateUser(modules.CreateUserRequest{
		Name:  "Jane",
		Email: "jane@example.com",
		Age:   25,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestCreateUser_ValidationError_MissingName(t *testing.T) {
	mock := &MockUserRepository{}

	uc := users.NewUserUsecase(mock, nil)
	_, err := uc.CreateUser(modules.CreateUserRequest{
		Email: "jane@example.com",
	})
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestCreateUser_ValidationError_MissingEmail(t *testing.T) {
	mock := &MockUserRepository{}

	uc := users.NewUserUsecase(mock, nil)
	_, err := uc.CreateUser(modules.CreateUserRequest{
		Name: "Jane",
	})
	if err == nil {
		t.Fatal("expected error for missing email")
	}
}

func TestCreateUser_PasswordHashed(t *testing.T) {
	var savedPassword string
	mock := &MockUserRepository{
		CreateUserFunc: func(req modules.CreateUserRequest) (int, error) {
			savedPassword = req.Password
			return 1, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	_, err := uc.CreateUser(modules.CreateUserRequest{
		Name:     "Jane",
		Email:    "jane@example.com",
		Password: "plaintext123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if savedPassword == "plaintext123" {
		t.Fatal("password should be hashed, not stored as plaintext")
	}
	if savedPassword == "" {
		t.Fatal("password should not be empty after hashing")
	}
}

func TestCreateUser_DefaultRole(t *testing.T) {
	var savedRole string
	mock := &MockUserRepository{
		CreateUserFunc: func(req modules.CreateUserRequest) (int, error) {
			savedRole = req.Role
			return 1, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	_, _ = uc.CreateUser(modules.CreateUserRequest{
		Name:  "Jane",
		Email: "jane@example.com",
	})
	if savedRole != "user" {
		t.Fatalf("expected default role 'user', got '%s'", savedRole)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mock := &MockUserRepository{
		UpdateUserFunc: func(id int, req modules.UpdateUserRequest) error {
			return nil
		},
	}

	name := "Updated"
	uc := users.NewUserUsecase(mock, nil)
	err := uc.UpdateUser(1, modules.UpdateUserRequest{Name: &name})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	mock := &MockUserRepository{
		UpdateUserFunc: func(id int, req modules.UpdateUserRequest) error {
			return fmt.Errorf("user with id %d not found or already deleted", id)
		},
	}

	name := "Updated"
	uc := users.NewUserUsecase(mock, nil)
	err := uc.UpdateUser(999, modules.UpdateUserRequest{Name: &name})
	if err == nil {
		t.Fatal("expected error for not found user")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mock := &MockUserRepository{
		DeleteUserByIDFunc: func(id int) (int64, error) {
			return 1, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	rows, err := uc.DeleteUserByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rows != 1 {
		t.Fatalf("expected 1 row affected, got %d", rows)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	mock := &MockUserRepository{
		DeleteUserByIDFunc: func(id int) (int64, error) {
			return 0, fmt.Errorf("user with id %d not found or already deleted", id)
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	_, err := uc.DeleteUserByID(999)
	if err == nil {
		t.Fatal("expected error for not found user")
	}
}

func TestRegisterUser_MissingPassword(t *testing.T) {
	mock := &MockUserRepository{}

	uc := users.NewUserUsecase(mock, nil)
	_, err := uc.RegisterUser(modules.CreateUserRequest{
		Name:  "Jane",
		Email: "jane@example.com",
	})
	if err == nil {
		t.Fatal("expected error for missing password")
	}
}

func TestLoginUser_InvalidEmail(t *testing.T) {
	mock := &MockUserRepository{
		GetUserByEmailFunc: func(email string) (*modules.User, error) {
			return nil, fmt.Errorf("user not found")
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	_, err := uc.LoginUser("nonexistent@example.com", "password")
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestCountUsers_Success(t *testing.T) {
	mock := &MockUserRepository{
		CountUsersFunc: func() (int, error) {
			return 42, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	count, err := uc.CountUsers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 42 {
		t.Fatalf("expected 42, got %d", count)
	}
}

func TestCreateUserWithAudit_Success(t *testing.T) {
	mock := &MockUserRepository{
		CreateUserWithAuditFunc: func(req modules.CreateUserRequest, action, details string) (int, error) {
			if action != "USER_CREATED" {
				return 0, fmt.Errorf("expected action USER_CREATED, got %s", action)
			}
			return 5, nil
		},
	}

	uc := users.NewUserUsecase(mock, nil)
	id, err := uc.CreateUserWithAudit(modules.CreateUserRequest{
		Name:  "Audit User",
		Email: "audit@example.com",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 5 {
		t.Fatalf("expected id 5, got %d", id)
	}
}
