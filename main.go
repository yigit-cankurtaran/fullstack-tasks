package main

// import http and gin
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// store task data in memory
type task struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Completion bool   `json:"completion"`
}

// seed task data
var tasks = []task{
	{ID: 1, Name: "You can create tasks", Completion: false},
	{ID: 2, Name: "You can read tasks", Completion: false},
	{ID: 3, Name: "You can update tasks", Completion: true},
	{ID: 4, Name: "You can delete tasks", Completion: false},
}

// main function
func main() {
	// if tasks.json exists, read from it
	tasksJSON, err := os.ReadFile("tasks.json")
	if err == nil {
		json.Unmarshal(tasksJSON, &tasks)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getTasks(w, r)
		case "POST":
			postTasks(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})

	mux.HandleFunc("/tasks/", taskByIDHandler)
	http.ListenAndServe(":8080", mux)

}

// GET /tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// POST /tasks
func postTasks(w http.ResponseWriter, r *http.Request) {
	var newTask task
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &newTask)
	// pointer because we want to modify the original slice
	if err != nil {
		http.Error(w, "error unmarshalling request body", http.StatusBadRequest)
		return
	}

	// add the new task to the slice
	tasks = append(tasks, newTask)

	// save to file
	tasksJSON, _ := json.Marshal(tasks)
	os.WriteFile("tasks.json", tasksJSON, os.ModePerm)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// GET /tasks/:id
func getTaskByID(c *gin.Context) {
	id := c.Param("id")
	// this is not about the "ID" variable in task struct, but the "id" variable in the URL
	// IMPORTANT!
	fmt.Println(id)

	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	// loop over the list of tasks, looking for
	// a task whose ID value matches the parameter
	// send a JSON response containing the task data
	for _, a := range tasks {
		if a.ID == idInt {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// PUT /tasks/:id
func putTaskByID(c *gin.Context) {
	id := c.Param("id")

	// get the task with the matching id
	var newTask task
	if err := c.BindJSON(&newTask); err != nil {
		return
	}

	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	for i, a := range tasks {
		if a.ID == idInt {
			tasks[i] = newTask // <-- update task data

			// save to file
			tasksJSON, err := json.Marshal(tasks)
			if err != nil {
				fmt.Println(err)
			}
			os.WriteFile("tasks.json", tasksJSON, os.ModePerm)

			c.IndentedJSON(http.StatusOK, newTask)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}

// DELETE /tasks/:id
func deleteTaskByID(c *gin.Context) {
	id := c.Param("id")

	// convert id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid ID"})
		return
	}

	for i, a := range tasks {
		if a.ID == idInt {
			tasks = append(tasks[:i], tasks[i+1:]...)

			// save to file
			tasksJSON, err := json.Marshal(tasks)
			if err != nil {
				fmt.Println(err)
			}
			os.WriteFile("tasks.json", tasksJSON, os.ModePerm)

			c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
}
