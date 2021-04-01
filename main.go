package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]

		fmt.Fprintf(w, "You've requested the book: %s on page %d\n", title, time.Now().UnixNano())
	})

	http.ListenAndServe(":80", r)
}
