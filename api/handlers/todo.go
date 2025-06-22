package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"github.com/Aergiaaa/todolist-on-vercel/models"
	"github.com/Aergiaaa/todolist-on-vercel/storage"
)

// TodoHandler handles todo-related requests
type TodoHandler struct {
	store     storage.TodoStorage
	templates *template.Template
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler(store storage.TodoStorage, templates *template.Template) *TodoHandler {
	return &TodoHandler{
		store:     store,
		templates: templates,
	}
}

// ListTodos handles GET requests to list all todos
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	var err error
	todos, err := h.store.GetAll()
	if err != nil {
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}

	// Handle HTMX request vs regular request
	if r.Header.Get("HX-Request") == "true" {
		err = h.templates.ExecuteTemplate(w, "todo-list", todos)
		if err != nil {
			http.Error(w, "Failed to render todo list", http.StatusInternalServerError)
			return
		}
	} else {
		err = h.templates.ExecuteTemplate(w, "index", todos)
		if err != nil {
			http.Error(w, "Failed to render index", http.StatusInternalServerError)
			return
		}
	}
}

// GetTodoForm handles GET requests to get the form for creating/editing a todo
func (h *TodoHandler) GetTodoForm(w http.ResponseWriter, r *http.Request) {

	// Instead of using a Todo struct, let's use a simple map for now to diagnose
	formData := map[string]string{
		"Title":       "",
		"Description": "",
	}

	id := r.URL.Query().Get("id")
	if id != "" {
		todo, err := h.store.Get(id)
		if err != nil {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		formData["Title"] = todo.Title
		formData["Description"] = todo.Description
		formData["ID"] = todo.ID
	} // Use the collected data when rendering the template	log.Println("Executing template 'todo-form', with data:", formData)
	err := h.templates.ExecuteTemplate(w, "todo-form", formData)
	if err != nil {
		log.Println("Error rendering form:", err)
		// Return a plain form as fallback
		fmt.Fprintf(w, "<div class='todo-form'><h2>Create New Todo</h2><form action='/todos' method='post'><div class='form-group'><label>Title:</label><input type='text' name='title'></div><div class='form-group'><label>Description:</label><textarea name='description'></textarea></div><button type='submit'>Create</button></form></div>")
		return
	}
}

// CreateTodo handles POST requests to create a new todo
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" {
		log.Println("Error: Title is required")
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo := models.NewTodo(title, description)

	err := h.store.Create(todo)
	if err != nil {
		log.Println("Error creating todo:", err)
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	todos, err := h.store.GetAll()
	if err != nil {
		log.Println("Error getting todos:", err)
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}

	// First clear the form container
	fmt.Fprintf(w, "<div id='form-container' hx-swap-oob='true'></div>")

	// Then render todo list
	err = h.templates.ExecuteTemplate(w, "todo-list", todos)
	if err != nil {
		log.Println("Error rendering todo-list:", err)
		http.Error(w, "Failed to render todo list", http.StatusInternalServerError)
		return
	}

	// Add the "Add Todo" button back with OOB swap
	fmt.Fprintf(w, "<section id='actions-container' class='actions' hx-swap-oob='true'>"+
		"<button class='add-btn' hx-get='/todos/form' hx-target='#form-container' hx-swap='innerHTML' "+
		"hx-swap-oob='true' hx-target='#actions-container' hx-swap-oob='outerHTML:#actions-container:none'>Add Todo</button></section>")

}

// UpdateTodo handles POST requests to update an existing todo
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		log.Println("Error: ID is required for update")
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	todo, err := h.store.Get(id)
	if err != nil {
		log.Println("Error getting todo with ID:", id, err)
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" {
		log.Println("Error: Title is required for update")
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	todo.Update(title, description)

	err = h.store.Update(todo)
	if err != nil {
		log.Println("Error updating todo in storage:", err)
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	todos, err := h.store.GetAll()
	if err != nil {
		log.Println("Error getting todos after update:", err)
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}
	// First clear the form container
	fmt.Fprintln(w, "<div id='form-container' hx-swap-oob='true'></div>")

	// Then render todo list
	err = h.templates.ExecuteTemplate(w, "todo-list", todos)
	if err != nil {
		http.Error(w, "Failed to render todo list", http.StatusInternalServerError)
		return
	}

	// Add the "Add Todo" button back with OOB swap
	fmt.Fprintln(w, "<section id='actions-container' class='actions' hx-swap-oob='true'>"+
		"<button class='add-btn' hx-get='/todos/form' hx-target='#form-container' hx-swap='innerHTML' "+
		"hx-swap-oob='true' hx-target='#actions-container' hx-swap-oob='outerHTML:#actions-container:none'>Add Todo</button></section>")
}

// ToggleTodoStatus handles POST requests to toggle a todo's completion status
func (h *TodoHandler) ToggleTodoStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	todo, err := h.store.Get(id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todo.ToggleStatus()
	err = h.store.Update(todo)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	err = h.templates.ExecuteTemplate(w, "todo-item", todo)
	if err != nil {
		http.Error(w, "Failed to render todo item", http.StatusInternalServerError)
		return
	}
}

// DeleteTodo handles DELETE requests to delete a todo
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	err := h.store.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ServeHTTP is the main entry point for the handler
func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log the request path to help with debugging

	// Handle the paths based on what's being stripped in main.go
	path := strings.TrimPrefix(r.URL.Path, "/")

	// Special handling for the root path or the /todos path
	if path == "" || path == "todos" {
		h.ListTodos(w, r)
		return
	}

	// For paths with /todos/something
	pathParts := strings.Split(path, "/")
	if len(pathParts) > 1 {
		subPath := pathParts[1]

		switch {
		case subPath == "form" && r.Method == http.MethodGet:
			h.GetTodoForm(w, r)
			return
		case subPath == "update" && r.Method == http.MethodPost:
			h.UpdateTodo(w, r)
			return
		case subPath == "toggle" && r.Method == http.MethodPost:
			h.ToggleTodoStatus(w, r)
			return
		case subPath == "delete" && r.Method == http.MethodDelete:
			h.DeleteTodo(w, r)
			return
		}
	}

	// For direct POST to /todos
	if r.Method == http.MethodPost {
		h.CreateTodo(w, r)
		return
	}
	// Default case
	http.NotFound(w, r)
}
