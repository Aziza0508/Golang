package users

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"golang/internal/repository"
	"golang/pkg/modules"
)

type Usecase struct {
	repo  repository.UserRepository
	redis *redis.Client
}

func NewUserUsecase(repo repository.UserRepository, redisClient *redis.Client) *Usecase {
	return &Usecase{
		repo:  repo,
		redis: redisClient,
	}
}

func (u *Usecase) GetUsers(limit, offset int) ([]modules.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return u.repo.GetUsers(limit, offset)
}

func (u *Usecase) GetUserByID(id int) (*modules.User, error) {
	if u.redis != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("user:%d", id)
		cached, err := u.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var user modules.User
			if json.Unmarshal([]byte(cached), &user) == nil {
				log.Printf("Cache HIT for user %d", id)
				return &user, nil
			}
		}
		log.Printf("Cache MISS for user %d", id)
	}

	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if u.redis != nil && user != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("user:%d", id)
		data, _ := json.Marshal(user)
		u.redis.Set(ctx, cacheKey, data, 5*time.Minute)
	}

	return user, nil
}

func (u *Usecase) CreateUser(req modules.CreateUserRequest) (int, error) {
	if req.Name == "" {
		return 0, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return 0, fmt.Errorf("email is required")
	}
	if req.Role == "" {
		req.Role = "user"
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return 0, fmt.Errorf("failed to hash password: %w", err)
		}
		req.Password = string(hashed)
	}
	return u.repo.CreateUser(req)
}

func (u *Usecase) UpdateUser(id int, req modules.UpdateUserRequest) error {
	err := u.repo.UpdateUser(id, req)
	if err != nil {
		return err
	}
	if u.redis != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("user:%d", id)
		u.redis.Del(ctx, cacheKey)
	}
	return nil
}

func (u *Usecase) DeleteUserByID(id int) (int64, error) {
	rowsAffected, err := u.repo.DeleteUserByID(id)
	if err != nil {
		return 0, err
	}
	if u.redis != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("user:%d", id)
		u.redis.Del(ctx, cacheKey)
	}
	return rowsAffected, nil
}

func (u *Usecase) CreateUserWithAudit(req modules.CreateUserRequest) (int, error) {
	if req.Name == "" {
		return 0, fmt.Errorf("name is required")
	}
	if req.Email == "" {
		return 0, fmt.Errorf("email is required")
	}
	if req.Role == "" {
		req.Role = "user"
	}
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return 0, fmt.Errorf("failed to hash password: %w", err)
		}
		req.Password = string(hashed)
	}
	return u.repo.CreateUserWithAudit(req, "USER_CREATED", fmt.Sprintf("User %s created with email %s", req.Name, req.Email))
}

func (u *Usecase) RegisterUser(req modules.CreateUserRequest) (int, error) {
	if req.Password == "" {
		return 0, fmt.Errorf("password is required")
	}
	return u.CreateUser(req)
}

func (u *Usecase) LoginUser(email, password string) (*modules.User, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	return user, nil
}

func (u *Usecase) CountUsers() (int, error) {
	return u.repo.CountUsers()
}
