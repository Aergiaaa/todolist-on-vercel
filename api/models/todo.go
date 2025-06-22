package models

import (
	"time"
)

// Todo represents a todo item
type Todo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewTodo creates a new todo item
func NewTodo(title, description string) *Todo {
	now := time.Now()
	return &Todo{
		ID:          generateID(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ToggleStatus toggles the completion status of the todo item
func (t *Todo) ToggleStatus() {
	t.Completed = !t.Completed
	t.UpdatedAt = time.Now()
}

// Update updates the todo item
func (t *Todo) Update(title, description string) {
	t.Title = title
	t.Description = description
	t.UpdatedAt = time.Now()
}

// generateID generates a simple ID for a todo item
// In a real application, you might want to use UUID or similar
func generateID() string {
	return time.Now().Format("20060102150405")
}
