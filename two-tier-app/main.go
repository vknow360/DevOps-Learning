package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
)

type Note struct {
	ID      int    `json:"id"`
	Message string `json:"message,omitempty"`
}

func main() {
	db := connectDB()
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Get all notes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		notes := getAllNotes(db)
		json.NewEncoder(w).Encode(notes)
	})

	// Add a note
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		note := addNote(db, r.FormValue("message"))
		if note.ID == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add note"})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	})

	// Get a note
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
			return
		}
		note := getNote(db, id)
		if note.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Note not found"})
			return
		}
		json.NewEncoder(w).Encode(note)
	})

	// Delete a note
	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID"})
			return
		}
		err = deleteNote(db, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete note"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	fmt.Println("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}

}

func connectDB() *sql.DB {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	pswd := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pswd, host, dbName)
	connStr = "root:root@tcp(localhost:3306)/notesdb"

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to database")
	fmt.Println("Running migrations...")

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS notes (id INT AUTO_INCREMENT PRIMARY KEY, message TEXT)")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Migrations complete")

	return db
}

func getAllNotes(db *sql.DB) []Note {
	query := `SELECT id FROM notes`
	rows, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		notes = append(notes, note)
	}

	return notes
}

func getNote(db *sql.DB, id int) Note {
	query := `SELECT * FROM notes WHERE id = ?`
	rows, err := db.Query(query, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	rows.Next()
	var note Note
	err = rows.Scan(&note.ID, &note.Message)
	if err != nil {
		fmt.Println(err.Error())
		return Note{}
	}
	fmt.Println(note)
	return note
}

func addNote(db *sql.DB, message string) Note {
	query := `INSERT INTO notes (message) VALUES (?)`
	result, err := db.Exec(query, message)
	if err != nil {
		panic(err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	return Note{ID: int(id)}
}

func deleteNote(db *sql.DB, id int) error {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
