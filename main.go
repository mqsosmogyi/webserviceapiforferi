package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	//"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	_ "github.com/go-sql-driver/mysql"
)

var router *chi.Mux
var db *sql.DB

type Item struct {
	Id   int    `json:"id"`
	Task string `json:"task`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func getTasks(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, task FROM user2")
	handleErr(err)
	defer rows.Close()

	var content []Item

	for rows.Next() {

		var Task Item
		var Id Item

		err := rows.Scan(&Id.Id, &Task.Task)
		handleErr(err)

		//save := fmt.Sprint(Id, Task)
		//fmt.Println(save)

		content = append(content, Task)

	}

	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, content)

	defer rows.Close()

}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func createTask(w http.ResponseWriter, r *http.Request) {

	var task Item

	json.NewDecoder(r.Body).Decode(&task)

	insert, err := db.Prepare("INSERT INTO user2(task) VALUES (?)")
	handleErr(err)

	_, er := insert.Exec(task.Task)
	handleErr(er)

	fmt.Println("Task Added!")

	defer insert.Close()

	name := task.Task
	w.Write([]byte("Task Added: "))
	w.Write([]byte(name))
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func deleteTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	delete, err := db.Prepare("DELETE FROM user2 WHERE id=?")
	handleErr(err)

	_, er := delete.Exec(id)
	handleErr(er)

	fmt.Println("Task deleted!")

	defer delete.Close()

	w.Write([]byte("Task Deleted!"))

}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func handleErr(err error) {

	if err != nil {
		panic(err)
	}
}

func init() {

	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	dbSource := fmt.Sprintf("root:GFDdsdf1234354231@tcp(127.0.0.1:3306)/golangdb")

	var err error
	db, err = sql.Open("mysql", dbSource)
	handleErr(err)
}

func main() {

	router.Delete("/test/{id}", deleteTask)
	router.Post("/test", createTask)
	router.Get("/", getTasks)

	fmt.Println("Listening...")
	http.ListenAndServe(":3000", router)
}
