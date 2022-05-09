package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
)

var router *chi.Mux
var db *sql.DB

type receipt struct {
	Id                 int    `json:"id"`
	Merchant_id        string `json:"merchant_id"`
	Owner_id           int    `json:"owner_id"`
	Total_amount       int    `json:"total_amount"`
	Payment_method     string `json:"payment_method"`
	Cassa_number       int    `json:"cassa_number"`
	Transaction_number int    `json:"transaction_number"`
	Ap_number          int    `json:"ap_number"`
	Created_at         int    `json:"created_at"`
}

type line_items struct {
	receipt
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	Measurement int    `json:"measurement"`
	Price       int    `json:"price"`
}

func main() {

	router.Get("/getreceipts/", getReceipts)
	router.Get("/getmerchant/{merchant_id}", getMerchantId)
	router.Get("/getownerid/{owner_id}", getOwnerId)
	router.Post("/createreceipts/", createReceipts)

	fmt.Println("Listening...")
	http.ListenAndServe(":3001", router)

}

func getMerchantId(w http.ResponseWriter, r *http.Request) {

	merchant_id := chi.URLParam(r, "merchant_id")

	rows, err := db.Query("SELECT * FROM receipt WHERE merchant_id = (?)", merchant_id)
	handleErr(err)
	defer rows.Close()

	content := []receipt{}

	for rows.Next() {

		cont := receipt{}

		err := rows.Scan(&cont.Id, &cont.Merchant_id, &cont.Owner_id, &cont.Total_amount, &cont.Payment_method, &cont.Cassa_number, &cont.Transaction_number, &cont.Ap_number, &cont.Created_at)
		handleErr(err)

		content = append(content, cont)
	}

	fmt.Println(content)

	defer rows.Close()

	data := content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)

}

func getOwnerId(w http.ResponseWriter, r *http.Request) {

	owner_id := chi.URLParam(r, "owner_id")

	rows, err := db.Query("SELECT * FROM receipt WHERE owner_id = (?)", owner_id)
	handleErr(err)
	defer rows.Close()

	content := []receipt{}

	for rows.Next() {

		cont := receipt{}

		err := rows.Scan(&cont.Id, &cont.Merchant_id, &cont.Owner_id, &cont.Total_amount, &cont.Payment_method, &cont.Cassa_number, &cont.Transaction_number, &cont.Ap_number, &cont.Created_at)
		handleErr(err)

		content = append(content, cont)

	}

	fmt.Println(content)

	defer rows.Close()

	data := content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)

}

func createReceipts(w http.ResponseWriter, r *http.Request) {

	var values receipt

	json.NewDecoder(r.Body).Decode(&values)

	fmt.Println(values)

	insert, err := db.Prepare("INSERT INTO receipt(merchant_id, owner_id, total_amount, payment_method, cassa_number, transaction_number, ap_number, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	handleErr(err)

	_, er := insert.Exec(values.Merchant_id, values.Owner_id, values.Total_amount, values.Payment_method, values.Cassa_number, values.Transaction_number, values.Ap_number, values.Created_at)
	handleErr(er)

	fmt.Println("Task Added!")

	defer insert.Close()
}

func getReceipts(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM receipt")
	handleErr(err)
	defer rows.Close()

	content := []receipt{}

	for rows.Next() {

		cont := receipt{}

		err := rows.Scan(&cont.Id, &cont.Merchant_id, &cont.Owner_id, &cont.Total_amount, &cont.Payment_method, &cont.Cassa_number, &cont.Transaction_number, &cont.Ap_number, &cont.Created_at)
		handleErr(err)

		content = append(content, cont)
	}

	fmt.Println(content)

	defer rows.Close()

	data := content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)

}

func init() {

	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	dbSource := fmt.Sprintf("root:GFDdsdf1234354231@tcp(127.0.0.1:3306)/golangdb")

	var err error
	db, err = sql.Open("mysql", dbSource)
	handleErr(err)
}

func handleErr(err error) {

	if err != nil {
		panic(err)
	}
}
