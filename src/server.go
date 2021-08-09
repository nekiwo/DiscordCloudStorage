package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func server() {
	// Static Pages
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Routing
	http.HandleFunc("/upload", func(res http.ResponseWriter, req *http.Request) {
		file, handler, err := req.FormFile("file")
		if err != nil {fmt.Println(err)}
		defer file.Close()

		req.MultipartReader()

		// Create a directory for the file
		FileID := strconv.FormatInt(time.Now().Unix(), 10);
		os.MkdirAll("temp/file" + FileID, os.ModePerm)

		// Create temp file
		TempFile, err := ioutil.TempFile("temp", "upload" + FileID + filepath.Ext(handler.Filename))
		if err != nil {fmt.Println(err)}
		defer TempFile.Close()

		// Copy over the file
		FileBytes, err := io.ReadAll(file)
		if err != nil {fmt.Println(err)}

		// Write file's data on the new temp file
		TempFile.Write(FileBytes)

		FileInfo, _ := TempFile.Stat()
		FileSize := FileInfo.Size()

		ChunkSize := 8 << 20
		TotalChunks := int(math.Ceil(float64(FileSize) / float64(ChunkSize)))

		for i := 0; i < TotalChunks; i++ {
			// Use regular size unless it's the last piece
			CurrentChunkSize := int(math.Min(float64(ChunkSize), float64(FileSize - int64(i * ChunkSize))))

			// Slice the data into a chunk
			ChunkData := make([]byte, CurrentChunkSize)
			file.Read(ChunkData)

			// Create file
			FileName := "./temp/file" + FileID + "/chunk" + strconv.Itoa(i)
			_, err := os.Create(FileName)
			if err != nil {fmt.Println(err); os.Exit(1)}

			// Save file
			ioutil.WriteFile(FileName, ChunkData, os.ModeAppend)
		}

		// Collect all metadata for future reference
		// Upload it to discord
		UploadFile(MetaData{
			handler.Filename,
			FileID,
			TotalChunks,
		})


		//https://socketloop.com/tutorials/golang-how-to-split-or-chunking-a-file-to-smaller-pieces
	})

	// Server stuff
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type MetaData struct {
	FileName string
	FileID string
	TotalChunks int
}