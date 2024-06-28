package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
	"os"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Task struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	ActiveAt time.Time `json:"activeAt"`
	Done     bool      `json:"done"`
}

var tasks = struct {
	sync.RWMutex
	m map[string]Task
}{m: make(map[string]Task)}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/todo-list/tasks", createTask).Methods("POST")
	r.HandleFunc("/api/todo-list/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/api/todo-list/tasks/{id}", deleteTask).Methods("DELETE")
	r.HandleFunc("/api/todo-list/tasks/{id}/done", markTaskDone).Methods("PUT")
	r.HandleFunc("/api/todo-list/tasks", listTasks).Methods("GET")
	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (t *Task) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		ActiveAt string `json:"activeAt"`
		Done     bool   `json:"done"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.ID = aux.ID
	t.Title = aux.Title
	t.Done = aux.Done

	parsedTime, err := time.Parse("2006-01-02", aux.ActiveAt)
	if err != nil {
		return err
	}
	t.ActiveAt = parsedTime

	return nil
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input: unable to parse JSON", http.StatusBadRequest)
		return
	}
	if len(t.Title) > 200 {
		http.Error(w, "Title too long", http.StatusBadRequest)
		return
	}
	if t.ActiveAt.IsZero() {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	t.ID = uuid.New().String()
	tasks.Lock()
	defer tasks.Unlock()
	for _, task := range tasks.m {
		if task.Title == t.Title && task.ActiveAt.Equal(t.ActiveAt) {
			http.Error(w, "Task already exists", http.StatusConflict)
			return
		}
	}
	tasks.m[t.ID] = t
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": t.ID})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if len(t.Title) > 200 {
		http.Error(w, "Title too long", http.StatusBadRequest)
		return
	}
	if t.ActiveAt.IsZero() {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}
	tasks.Lock()
	defer tasks.Unlock()
	task, exists := tasks.m[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	t.ID = task.ID
	t.Done = task.Done
	tasks.m[id] = t
	w.WriteHeader(http.StatusNoContent)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tasks.Lock()
	defer tasks.Unlock()
	if _, exists := tasks.m[id]; !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	delete(tasks.m, id)
	w.WriteHeader(http.StatusNoContent)
}

func markTaskDone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tasks.Lock()
	defer tasks.Unlock()
	task, exists := tasks.m[id]
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	task.Done = true
	tasks.m[id] = task
	w.WriteHeader(http.StatusNoContent)
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	tasks.RLock()
	defer tasks.RUnlock()
	var result []Task
	for _, task := range tasks.m {
		if status == "done" && task.Done {
			result = append(result, task)
		} else if status != "done" && !task.Done && task.ActiveAt.Before(time.Now()) {
			if task.ActiveAt.Weekday() == time.Saturday || task.ActiveAt.Weekday() == time.Sunday {
				task.Title = "ВЫХОДНОЙ - " + task.Title
			}
			result = append(result, task)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
