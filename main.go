package main

import (
	"html/template"
	"net/http"
)

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/SendRequest", RequestHandler)
	http.ListenAndServe(":3030", nil)
}

func viewIndex(w *http.ResponseWriter) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		wr := *w
		wr.Write([]byte(err.Error()))
	}

	t.ExecuteTemplate(*w, "index", nil)
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	viewIndex(&w)
}

func RequestHandler(w http.ResponseWriter, req *http.Request) {
	request := req.FormValue("request")
	viewIndex(&w)
	w.Write([]byte("<h1>Results for " + request + ":</h1>"))
	for i := 0; i <= 20; i++ {
		w.Write([]byte("<a href=" + "http://amazon.com>" + "Amazon<a><br>"))
	}
	http.Redirect(w, req, "/", 0)
}
