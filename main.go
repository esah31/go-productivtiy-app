package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

// Task represents a to-do task
type Task struct {
    ID        int       `json:"id"`
    Text      string    `json:"text"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"createdAt"`
}

var db *sql.DB

func main() {
    // Initialize database
    initDB()
    defer db.Close()

    // Set up routes
    r := mux.NewRouter()
    r.HandleFunc("/api/tasks", getTasksHandler).Methods("GET")
    r.HandleFunc("/api/tasks", addTaskHandler).Methods("POST")
    r.HandleFunc("/api/tasks/{id}", deleteTaskHandler).Methods("DELETE")
    r.HandleFunc("/api/tasks/{id}/complete", completeTaskHandler).Methods("PUT")

    // Serve static files
    staticDir := http.Dir("./static/")
    staticFileHandler := http.FileServer(staticDir)
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileHandler))

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Println("Server running on port", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func initDB() {
    var err error
    db, err = sql.Open("sqlite3", "./database.db")
    if err != nil {
        log.Fatal("Error opening database:", err)
    }

    // Create tasks table if it doesn't exist
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL,
        completed BOOLEAN NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`)
    if err != nil {
        log.Fatal("Error creating tasks table:", err)
    }
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, text, completed, created_at FROM tasks ORDER BY created_at DESC")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var tasks []Task
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Text, &task.Completed, &task.CreatedAt)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        tasks = append(tasks, task)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
    var task Task
    err := json.NewDecoder(r.Body).Decode(&task)
    if err != nil || task.Text == "" {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("INSERT INTO tasks (text, completed) VALUES (?, ?)", task.Text, false)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    taskID := params["id"]

    _, err := db.Exec("DELETE FROM tasks WHERE id = ?", taskID)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    taskID := params["id"]

    _, err := db.Exec("UPDATE tasks SET completed = ? WHERE id = ?", true, taskID)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}
