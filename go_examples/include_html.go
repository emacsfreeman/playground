package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const (
	Port = ":8080"
)

func serveStatic(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("test.html")
	if err != nil {
		fmt.Println(err)
	}
	items := struct {
		Country string
		City    string
	}{
		Country: "Australia",
		City:    "Paris",
	}
	t.Execute(w, items)
}

func main() {
	http.HandleFunc("/", serveStatic)
	http.ListenAndServe(Port, nil)
}
