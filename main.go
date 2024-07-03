package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type task struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Completion bool   `json:"completion"`
}

var tasks = []task{
	{ID: 1, Name: "You can create tasks", Completion: false},
	{ID: 2, Name: "You can read tasks", Completion: false},
	{ID: 3, Name: "You can update tasks", Completion: true},
	{ID: 4, Name: "You can delete tasks", Completion: false},
}

func main() {
	tasksJSON, err := os.ReadFile("tasks.json")
	if err != nil {
		log.Println("Error reading tasks.json:", err)
	} else {
		if err := json.Unmarshal(tasksJSON, &tasks); err != nil {
			log.Println("Error unmarshalling tasks.json:", err)
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Accessed /tasks with method:", r.Method)
		if r.URL.Path != "/tasks" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case "GET":
			getTasks(w, r)
		case "POST":
			postTasks(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Accessed /tasks/ with method:", r.Method)
		taskByIDHandler(w, r)
	})

	log.Println("Server is running on localhost:1239")
	if err := http.ListenAndServe("localhost:1239", mux); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Println("Error encoding tasks:", err)
	}
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	var newTask task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		log.Println("Error reading request body:", err)
		return
	}
	if err := json.Unmarshal(body, &newTask); err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		log.Println("Error unmarshalling request body:", err)
		return
	}
	tasks = append(tasks, newTask)
	if tasksJSON, err := json.Marshal(tasks); err == nil {
		os.WriteFile("tasks.json", tasksJSON, os.ModePerm)
	} else {
		log.Println("Error marshalling tasks:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		log.Println("Error encoding new task:", err)
	}
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		log.Println("Invalid ID:", idStr)
		return
	}

	switch r.Method {
	case "GET":
		for _, t := range tasks {
			if t.ID == id {
				if err := json.NewEncoder(w).Encode(t); err != nil {
					log.Println("Error encoding task:", err)
				}
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)
		log.Println("Task not found with ID:", id)
	// Implement PUT and DELETE as needed
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
