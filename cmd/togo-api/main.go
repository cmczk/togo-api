package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cmczk/todo-api/internal/database"
	"github.com/cmczk/todo-api/internal/handler"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	log.Println("connecting to database...")
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err.Error())
	}
	defer db.Close()
	log.Println("connection to database has been established")

	store := database.NewTodoStore(db)
	handlers := handler.NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("/todos", methodHandler(handlers.GetAllTodos, http.MethodGet))
	mux.HandleFunc("/todos/new", methodHandler(handlers.CreateTodo, http.MethodPost))
	mux.HandleFunc("/todos/", todoIDHandler(handlers))

	loggingMiddleware := handler.LoggingMiddleware(mux)

	if err := http.ListenAndServe(":"+port, loggingMiddleware); err != nil {
		log.Fatalf("cannot run server: %s", err.Error())
	}
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}

		handlerFunc(w, r)
	}
}

func todoIDHandler(handler *handler.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTodoByID(w, r)
		case http.MethodPut:
			handler.UpdateTodo(w, r)
		case http.MethodDelete:
			handler.DeleteTodo(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
