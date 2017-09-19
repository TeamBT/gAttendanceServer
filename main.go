package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://jwuzvvbakkswqj:638c8277583ccec74041d91560bd319728949fe1a8085eb0570ad9e7e29cbf3c@ec2-54-243-255-57.compute-1.amazonaws.com:5432/derfn0d9nq69b3")

	if err != nil {
		panic(err)
	}
}

//Student ...
type Student struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Rfid      string `json:"rfid"`
	Partial   string `json:"partial"`
	CheckedIn bool   `json:"checkedIn"`
	Excused   bool   `json:"excused"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", redirectStudent)
	http.HandleFunc("/student", studentsIndex)
	http.HandleFunc("/student/show", studentShow)
	// http.HandleFunc("/students/create", createStudent)
	http.HandleFunc("/student/update", updateStudent)
	http.HandleFunc("/student/delete", deleteStudent)
	http.HandleFunc("/student/reset", resetStudents)

	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func redirectStudent(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/student", http.StatusSeeOther)
}

func studentsIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM student")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	studs := make([]Student, 0)
	for rows.Next() {

		stud := Student{}
		err = rows.Scan(&stud.ID, &stud.Name, &stud.Rfid, &stud.Partial, &stud.CheckedIn, &stud.Excused)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		studs = append(studs, stud)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	js, err := json.Marshal(studs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func studentShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM student WHERE id = $1", id)

	stud := Student{}
	err := row.Scan(&stud.ID, &stud.Name, &stud.Rfid, &stud.Partial, &stud.CheckedIn, &stud.Excused)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(stud)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

// func createStudent(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
// 		return
// 	}
//
// 	stud := Student{}
// 	stud.ID = r.FormValue("id")
// 	stud.Rfid = r.FormValue("rfid")
// 	stud.CheckedIn = r.FormValue("checkedIn") == "true"
// 	stud.Excused = r.FormValue("excused") == "true"
// }

func updateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	stud := Student{}
	stud.ID = r.FormValue("id")
	stud.Rfid = r.FormValue("rfid")
	stud.CheckedIn = r.FormValue("checkedIn") == "true"
	stud.Excused = r.FormValue("excused") == "true"

	if stud.ID != "" && stud.Rfid == "" {
		_, err := db.Exec("UPDATE student SET checked_in=$2, excused=$3 WHERE id=$1;", stud.ID, stud.CheckedIn, stud.Excused)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}

	if stud.Rfid != "" && stud.ID == "" {
		_, err := db.Exec("UPDATE student SET checked_in=$2 WHERE rfid=$1;", stud.Rfid, true)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
	}
}

func resetStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	_, err := db.Exec("UPDATE student SET checked_in=$1, excused=$2;", false, false)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/student", http.StatusSeeOther)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM student WHERE id=$1;", id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student", http.StatusSeeOther)
}
