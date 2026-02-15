package storage

import (
	"practice2/internal/models"
	"sort"
	"sync"
)

type TaskStorage struct {
	tasks  map[int]models.Task
	nextID int
	mu     sync.RWMutex
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *TaskStorage) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[s.nextID] = task
	s.nextID++
	return task
}

func (s *TaskStorage) GetByID(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	return task, exists
}

func (s *TaskStorage) GetAll() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result
}

func (s *TaskStorage) Update(id int, done bool) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, false
	}

	task.Done = done
	s.tasks[id] = task
	return task, true
}
