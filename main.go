package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	_ "github.com/go-sql-driver/mysql"
)

var (
	router *chi.Mux
	db     *sql.DB
)

type Item struct {
	Id   int    `json:"id"`
	Task string `json:"task`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func getTasks(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, task FROM user2")
	handleErr(err)
	defer rows.Close()

	content := []Item{}

	for rows.Next() {

		cont := Item{}

		err := rows.Scan(&cont.Id, &cont.Task)
		handleErr(err)

		content = append(content, cont)
	}

	new := []Item{}

	for i := len(content) - 1; i >= 0; i-- {
		new = append(new, content[i])
	}

	fmt.Println(new)
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, new)

	defer rows.Close()

}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func createTask(w http.ResponseWriter, r *http.Request) {

	var task Item

	task.Task = r.FormValue("nameName")

	json.NewDecoder(r.Body).Decode(&task.Task)

	fmt.Println(task.Task)

	insert, err := db.Prepare("INSERT INTO user2(task) VALUES (?)")
	handleErr(err)

	_, er := insert.Exec(task.Task)
	handleErr(er)

	http.Redirect(w, r, "https://f52c-178-48-68-68.eu.ngrok.io/", http.StatusSeeOther)

	fmt.Println("Task Added!")

	defer insert.Close()
}

//////////////////////////////////////////////////////////////////////////////////////////////////////

func deleteTask(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.FormValue("idtask")
	fmt.Print(id)

	delete, err := db.Prepare("DELETE FROM user2 WHERE (`id` = ?);")
	handleErr(err)

	query, er := delete.Exec(id)
	fmt.Println(query)
	handleErr(er)

	defer delete.Close()

	http.Redirect(w, r, "https://f52c-178-48-68-68.eu.ngrok.io/", http.StatusSeeOther)

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
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Delete("/delete/", deleteTask)
	router.Post("/test", createTask)
	router.Get("/", getTasks)

	fmt.Println("Listening...")
	http.ListenAndServe(":3000", router)
}
