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
	Receipt_id         int    `json:"receipt_id"`
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Quantity           int    `json:"quantity"`
	Measurement        int    `json:"measurement"`
	Price              int    `json:"price"`
}

/*

http://localhost:3001/getreceipts/ 			-> Get all of the receipts

http://localhost:3001/getmerchant/{merchant_id} 	-> http://localhost:3001/getmerchant/tesco -> Get merchant (tesco/penny..)

http://localhost:3001/getownerid/{owner_id} 		-> http://localhost:3001/getownerid/1 -> Get owner by ID

http://localhost:3001/getreceiptid/{id} 		-> http://localhost:3001/getreceiptid/2 -> Get receipt by ID

http://localhost:3001/createreceipts/ 			-> Add new receipt

    {
        "merchant_id": "lidl",
        "owner_id": 2,
        "total_amount": 200,
        "payment_method": "cash",
        "cassa_number": 3,
        "transaction_number": 1,
        "ap_number": 1,
        "created_at": 20220506
    }

http://localhost:3001/addlineitem/ -> Add new line item

    {
        "name": "sajt",
        "quantity": 2,
        "measurement": 15,
        "price": 5000,
        "receipt_id": 2
    }

*/

func main() {

	router.Get("/getreceipts/", getReceipts)
	router.Get("/getmerchant/{merchant_id}", getMerchantId)
	router.Get("/getownerid/{owner_id}", getOwnerId)
	router.Get("/getreceiptid/{id}", getReceiptId)
	router.Post("/createreceipts/", createReceipts)
	router.Post("/addlineitem/", addLineItem)

	fmt.Println("Listening...")
	http.ListenAndServe(":3001", router)

}

func getReceiptId(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	rows, err := db.Query("SELECT * FROM receipt WHERE id = (?)", id)
	handleErr(err)
	defer rows.Close()

	item_rows, er := db.Query("SELECT name, quantity, measurement, price, receipt_id FROM line_item WHERE receipt_id = (?)", id)
	handleErr(er)
	defer item_rows.Close()

	content := []receipt{}

	for rows.Next() {

		cont := receipt{}

		err := rows.Scan(&cont.Id, &cont.Merchant_id, &cont.Owner_id, &cont.Total_amount, &cont.Payment_method, &cont.Cassa_number, &cont.Transaction_number, &cont.Ap_number, &cont.Created_at)
		handleErr(err)

		content = append(content, cont)
	}
	defer rows.Close()

	items := []line_items{}

	for item_rows.Next() {

		item_list := line_items{}

		err := item_rows.Scan(&item_list.Name, &item_list.Quantity, &item_list.Measurement, &item_list.Price, &item_list.Receipt_id)
		handleErr(err)

		items = append(items, item_list)
	}
	defer item_rows.Close()

	item_data := items
	data := content

	headerResponse(w, r)

	json.NewEncoder(w).Encode(data)
	json.NewEncoder(w).Encode(item_data)
}

func addLineItem(w http.ResponseWriter, r *http.Request) {

	var items line_items

	json.NewDecoder(r.Body).Decode(&items)

	insert, err := db.Prepare("INSERT INTO line_item(name, quantity, measurement, price, receipt_id) VALUES (?, ?, ?, ?, ?)")
	handleErr(err)

	_, er := insert.Exec(&items.Name, &items.Quantity, &items.Measurement, &items.Price, &items.Receipt_id)
	handleErr(er)

	defer insert.Close()
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
	defer rows.Close()

	data := content
	headerResponse(w, r)
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
	defer rows.Close()

	data := content
	headerResponse(w, r)
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
	defer rows.Close()

	data := content
	headerResponse(w, r)
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

func headerResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func handleErr(err error) {

	if err != nil {
		panic(err)

	}
}
