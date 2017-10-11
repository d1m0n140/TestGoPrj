package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

type Table_view struct {
	id_client    string
	client_name  string
	client_phone string
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/SendRequest", RequestHandler)
	http.ListenAndServe(":3030", nil)
}

func viewIndexPage(w *http.ResponseWriter) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		wr := *w
		wr.Write([]byte(err.Error()))
	}

	t.ExecuteTemplate(*w, "index", nil)
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	viewIndexPage(&w)
}

func RequestHandler(w http.ResponseWriter, req *http.Request) {
	request := req.FormValue("request")

	var qResult string
	db, conEr := sql.Open("mssql", "server=localhost\\SQLExpress;user id=sa;password=mypassword;database=Example;")
	if conEr != nil {
		fmt.Printf("Connection error!!!: %s", conEr.Error())
	}

	rows, qEr := db.Query("SELECT * FROM Clients WHERE client_name=?;", request)

	if qEr != nil {
		fmt.Printf("An error occured while executing SQL-query: %s", qEr.Error())
	}

	tv := make([]*Table_view, 0)
	for rows.Next() {
		rw := new(Table_view)
		rows.Scan(&rw.id_client, &rw.client_name, &rw.client_phone)
		tv = append(tv, rw)
	}
	for _, rw := range tv {
		qResult += (rw.id_client + rw.client_name + rw.client_phone + "<br>")
	}

	defer rows.Close()
	defer db.Close()

	viewIndexPage(&w)
	w.Write([]byte("<h1>Results for " + request + ":</h1>"))
	for i := 0; i <= 20; i++ {
		w.Write([]byte(qResult))
	}
}
