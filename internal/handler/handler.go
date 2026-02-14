package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cmczk/todo-api/internal/database"
	"github.com/cmczk/todo-api/internal/models"
)

type Handler struct {
	store *database.TodoStore
}

func NewHandler(store *database.TodoStore) *Handler {
	return &Handler{
		store: store,
	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"message": message})
}

func getIDFromURLPath(urlPath string, prefix string) (int, error) {
	idStr := strings.Split(strings.TrimPrefix(urlPath, prefix), "/")[0]
	return strconv.Atoi(idStr)
}

func (h *Handler) getAllTodos(w http.ResponseWriter, _ *http.Request) {
	todos, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot get all todos")
		return
	}

	respondWithJSON(w, http.StatusOK, todos)
}

func (h *Handler) getTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURLPath(r.URL.Path, "/todos/")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	todo, err := h.store.GetByID(id)
	if err != nil {
		// TODO: add 404 case
		respondWithError(w, http.StatusInternalServerError, "cannot get todo by id")
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "cannot read request body")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "title cannot be empty")
		return
	}

	todo, err := h.store.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot create new todo")
		return
	}

	respondWithJSON(w, http.StatusCreated, todo)
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURLPath(r.URL.Path, "/todos/")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input models.UpdateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "cannot read request body")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "title cannot be empty")
		return
	}

	todo, err := h.store.Update(id, input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot create new todo")
		return
	}

	respondWithJSON(w, http.StatusOK, todo)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromURLPath(r.URL.Path, "/todos/")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.store.Delete(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "cannot delete todo")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
