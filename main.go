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
	order_date   string
	client_order string
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

	db, conEr := sql.Open("mssql", "server=localhost\\SQLExpress;user id=sa;password=mypassword;database=Example;")
	if conEr != nil {
		fmt.Printf("Connection error!!!: %s", conEr.Error())
	}

	query := "SELECT Clients.id_client, Clients.client_name, Clients.client_phone, order_date, Goods.name FROM Orders" +
		" INNER JOIN Clients ON Orders.id_client = Clients.id_client" +
		" INNER JOIN Goods ON Orders.id_good = Goods.id_good" +
		" WHERE client_name='" + request + "';"

	rows, qEr := db.Query(query)

	if qEr != nil {
		fmt.Printf("An error occured while executing SQL-query: %s", qEr.Error())
	}

	tv := make([]*Table_view, 0)
	for rows.Next() {
		rw := new(Table_view)
		rows.Scan(&rw.id_client, &rw.client_name, &rw.client_phone, &rw.order_date, &rw.client_order)
		tv = append(tv, rw)
	}

	defer rows.Close()
	defer db.Close()

	viewIndexPage(&w)

	w.Write([]byte("<h1>Results for " + request + ":</h1>"))
	w.Write([]byte(viewResult(tv)))
}

func viewResult(tv []*Table_view) string {
	qResult := "<table border=1 width=100% cellpadding=5>" +
		"<tr>" +
		"<th>Id</th>" +
		"<th>Name</th>" +
		"<th>Phone</th>" +
		"<th>Date</th>" +
		"<th>Ordered</th>" +
		"</tr>" +
		"<tr>"

	for _, rw := range tv {
		qResult += ("<td>" + rw.id_client + "</td>" + "<td>" + rw.client_name + "</td>" +
			"<td>" + rw.client_phone + "</td>" + "<td>" + rw.order_date + "</td>" +
			"<td>" + rw.client_order + "</td></tr>")
	}
	qResult += "</table>"
	return qResult
}
