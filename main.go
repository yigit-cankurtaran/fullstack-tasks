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
	// reading from tasks.json
	tasksJSON, err := os.ReadFile("tasks.json")
	if err != nil {
		if os.IsNotExist(err) {
			// Create the tasks.json file if it doesn't exist
			_, createErr := os.Create("tasks.json")
			if createErr != nil {
				log.Println("Error creating tasks.json:", createErr)
			}
		} else {
			log.Println("Error reading tasks.json:", err)
		}
	} else {
		if err := json.Unmarshal(tasksJSON, &tasks); err != nil {
			log.Println("Error unmarshalling tasks.json:", err)
		}
	}

	// create mux
	mux := http.NewServeMux()

	// handler for /tasks
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Accessed /tasks with method:", r.Method)
		// if not /tasks, return 404
		if r.URL.Path != "/tasks" {
			http.NotFound(w, r)
			return
		}
		// different functions based on HTTP method
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
	if err := http.ListenAndServe("localhost:1239", corsMiddleware(mux)); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	// set the header to application/json
	w.Header().Set("Content-Type", "application/json")
	// encode the tasks slice to JSON and write it to the response
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Println("Error encoding tasks:", err)
	}
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	// new var of type task
	var newTask task
	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		log.Println("Error reading request body:", err)
		return
	}
	// unmarshal the request body to newTask
	// unmarshaling = converting JSON to Go struct
	if err := json.Unmarshal(body, &newTask); err != nil {
		http.Error(w, "Error unmarshalling request body", http.StatusBadRequest)
		log.Println("Error unmarshalling request body:", err)
		return
	}
	// if unmarshaling is successful, append newTask to tasks
	tasks = append(tasks, newTask)
	// write the updated tasks slice to tasks.json
	if tasksJSON, err := json.Marshal(tasks); err == nil {
		os.WriteFile("tasks.json", tasksJSON, os.ModePerm)
	} else {
		log.Println("Error marshalling tasks:", err)
	}
	// set header to application/json and status to 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		log.Println("Error encoding new task:", err)
	}
}

func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// extract the ID from the URL
	idStr := r.URL.Path[len("/tasks/"):]
	// convert the ID to an integer
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

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// if the request method is OPTIONS, return immediately
		if r.Method == "OPTIONS" {
			return
		}

		// call the handler
		handler.ServeHTTP(w, r)
	})
}
