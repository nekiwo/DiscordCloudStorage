package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func server() {
	// Static Pages
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Routing
	http.HandleFunc("/upload", func(res http.ResponseWriter, req *http.Request) {
		req.MultipartReader()

		file, handler, err := req.FormFile("file")
		if err != nil {fmt.Println(err)}
		defer file.Close()

		// Create file
		dst, err := os.Create("./bin/" + handler.Filename)
		defer dst.Close()
		if err != nil {http.Error(res, err.Error(), http.StatusInternalServerError)}

		// Copy over the file
		if _, err := io.Copy(dst, file); err != nil {http.Error(res, err.Error(), http.StatusInternalServerError)}

		//https://socketloop.com/tutorials/golang-how-to-split-or-chunking-a-file-to-smaller-pieces
	})

	// Server stuff
	log.Fatal(http.ListenAndServe(":3000", nil))
}