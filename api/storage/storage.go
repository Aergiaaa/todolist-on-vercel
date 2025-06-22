package storage

import (
	"github.com/Aergiaaa/todolist-on-vercel/models"
)

// TodoStorage defines the interface for todo storage
type TodoStorage interface {
	GetAll() ([]*models.Todo, error)
	Get(id string) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id string) error
}
