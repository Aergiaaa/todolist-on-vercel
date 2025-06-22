package storage

import (
	"errors"
	"sync"

	"github.com/Aergiaaa/todolist-on-vercel/models"
)

// MemoryStorage implements TodoStorage using in-memory storage
type MemoryStorage struct {
	todos map[string]*models.Todo
	mutex sync.RWMutex
}

// NewMemoryStorage creates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		todos: make(map[string]*models.Todo),
	}
}

// GetAll returns all todo items
func (s *MemoryStorage) GetAll() ([]*models.Todo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todos := make([]*models.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

// Get returns a todo item by ID
func (s *MemoryStorage) Get(id string) (*models.Todo, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

// Create adds a new todo item
func (s *MemoryStorage) Create(todo *models.Todo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.todos[todo.ID] = todo
	return nil
}

// Update updates an existing todo item
func (s *MemoryStorage) Update(todo *models.Todo) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.todos[todo.ID]; !exists {
		return errors.New("todo not found")
	}
	s.todos[todo.ID] = todo
	return nil
}

// Delete removes a todo item
func (s *MemoryStorage) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.todos[id]; !exists {
		return errors.New("todo not found")
	}
	delete(s.todos, id)
	return nil
}
