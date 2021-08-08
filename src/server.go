package main

import (
	"log"
	"net/http"
)

func server() {
	// Static Pages
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Routing
	http.HandleFunc("/upload", upload)

	// Server stuff
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func upload(res http.ResponseWriter, req *http.Request) {

}