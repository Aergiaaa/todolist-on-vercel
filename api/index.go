package api

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Aergiaaa/todolist-on-vercel/handlers"
	"github.com/Aergiaaa/todolist-on-vercel/storage"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize the storage
	todoStore := storage.NewMemoryStorage()

	// Load templates
	templates := loadTemplates()

	// Initialize the handlers
	todoHandler := handlers.NewTodoHandler(todoStore, templates) // Set up routes

	mux := http.NewServeMux()
	mux.Handle("/todos/", todoHandler)
	mux.Handle("/todos", todoHandler)

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.NotFound(w, r)
			return
		}
		todoHandler.ListTodos(w, r)
	})

	mux.ServeHTTP(w, r)
}

// loadTemplates loads the HTML templates
func loadTemplates() *template.Template {
	// Create a new template with empty name
	tmpl := template.New("")

	// Get template files
	templateFiles, err := filepath.Glob("templates/*.html")
	if err != nil {
		log.Println("Warning: Failed to get template files:", err)
	}

	// debug
	for _, file := range templateFiles {
		log.Println("Template file found:", file)
	}

	// Parse templates
	tmpl, err = tmpl.ParseFiles(templateFiles...)
	if err != nil {
		log.Println("Warning: Failed to parse templates:", err)
		return tmpl
	}

	return tmpl
}
