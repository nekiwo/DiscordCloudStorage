package main

import (
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func server() {
	// Static Pages
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// Routing
	http.HandleFunc("/upload", func(res http.ResponseWriter, req *http.Request) {
		file, handler, err := req.FormFile("file")
		ErrCheck(err)
		defer file.Close()

		req.MultipartReader()

		// Create a directory for the file
		FileID := strconv.FormatInt(time.Now().Unix(), 10);
		os.MkdirAll("temp/file" + FileID, os.ModePerm)

		// Create temp file
		TempFile, err := ioutil.TempFile("temp", "upload*")
		ErrCheck(err)
		defer TempFile.Close()

		// Copy over the file
		FileBytes, err := io.ReadAll(file)
		ErrCheck(err)

		// Write file's data on the new temp file
		TempFile.Write(FileBytes)

		FileInfo, _ := TempFile.Stat()
		FileSize := FileInfo.Size()

		ChunkSize := 7 << 20
		TotalChunks := int(math.Ceil(float64(FileSize) / float64(ChunkSize)))

		for i := 0; i < TotalChunks; i++ {
			// Use regular size unless it's the last piece
			CurrentChunkSize := int(math.Min(float64(ChunkSize), float64(FileSize - int64(i * ChunkSize))))

			// Slice the data into a chunk
			ChunkData := make([]byte, CurrentChunkSize)
			file.Read(ChunkData)

			// Create file
			FileName := "./temp/file" + FileID + "/chunk" + strconv.Itoa(i)
			_, _ = os.Create(FileName)

			// Save file
			ioutil.WriteFile(FileName, ChunkData, os.ModeAppend)
		}

		// Delete temp file
		err = os.Remove("temp/" + TempFile.Name())
		ErrCheck(err)

		// Collect all metadata for future reference
		// Upload it to discord
		UploadFiles(MetaData{
			handler.Filename,
			FileID,
			TotalChunks,
		})


		//https://socketloop.com/tutorials/golang-how-to-split-or-chunking-a-file-to-smaller-pieces
	})

	http.HandleFunc("/download", func(res http.ResponseWriter, req *http.Request) {
		id, err := io.ReadAll(req.Body)
		ErrCheck(err)

		DownloadFiles(string(id))
	})

	// Server stuff
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type MetaData struct {
	FileName string
	FileID string
	TotalChunks int
}