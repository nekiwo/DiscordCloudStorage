package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

var TempFileList []TempFile

func server() {
	// Static Pages
	http.Handle("/", http.FileServer(http.Dir("public")))

	// Routing
	http.HandleFunc("/upload", func(res http.ResponseWriter, req *http.Request) {
		// Clean every file older than 1 day
		go cleanup()

		file, handler, err := req.FormFile("file")
		ErrCheck(err)
		defer file.Close()

		req.MultipartReader()

		// Create a directory for the file
		FileID := strconv.FormatInt(time.Now().Unix(), 10);
		os.MkdirAll("temp/file" + FileID, os.ModePerm)

		// Create temp file
		out, err := os.Create("temp/upload" + FileID)
		ErrCheck(err)
		defer out.Close()

		// Write file's data on the new temp file
		_, err = io.Copy(out, file)
		ErrCheck(err)

		FileBytes, err := ioutil.ReadFile(out.Name())
		ErrCheck(err)

		FileInfo, _ := out.Stat()
		FileSize := FileInfo.Size()
		ChunkSize := 7 << 20

		TotalChunks := int(math.Ceil(float64(FileSize) / float64(ChunkSize)))

		for i := 0; i < TotalChunks; i++ {
			// Use regular size unless it's the last piece
			CurrentSize := int(math.Min(float64(ChunkSize), float64(int(FileSize) - i * ChunkSize)))

			// Slice the data into a chunk
			ChunkData := FileBytes[i * ChunkSize: i * ChunkSize + CurrentSize]

			// Create file
			FileName := "temp/file" + FileID + "/chunk" + strconv.Itoa(i)
			chunk, err := os.Create(FileName)
			ErrCheck(err)
			defer chunk.Close()

			// Save file
			err = ioutil.WriteFile(FileName, ChunkData, os.ModeAppend)
			ErrCheck(err)
		}

		fmt.Println("Done generating chunks for " + FileID)

		// Collect all metadata for future reference
		// Upload it to discord
		UploadFiles(MetaData{
			handler.Filename,
			FileID,
			TotalChunks,
		})
	})

	http.HandleFunc("/download", func(res http.ResponseWriter, req *http.Request) {
		id, err := io.ReadAll(req.Body)
		ErrCheck(err)

		result := DownloadFiles(string(id))
		_, err = res.Write([]byte(result))
		ErrCheck(err)

		fmt.Println(result)
	})

	// Server stuff
	log.Fatal(http.ListenAndServe(":3000", nil))
}

type MetaData struct {
	FileName string
	FileID string
	TotalChunks int
}

type TempFile struct {
	FileName string
	FileID string
}

func cleanup() {
	for i := 0; i < len(TempFileList); i++ {
		IntTime, err := strconv.Atoi(TempFileList[i].FileID)
		ErrCheck(err)

		fmt.Println(IntTime)
		fmt.Println(int(time.Now().Unix()))
		fmt.Println(int(time.Now().Unix()) - IntTime)

		if int(time.Now().Unix()) - IntTime > 86400 {
			// Delete files if it's older than 24 hours
			err := os.Remove("public/files/u" + TempFileList[i].FileID + TempFileList[i].FileName)
			ErrCheck(err)

			err = os.Remove("temp/upload" + TempFileList[i].FileID)
			ErrCheck(err)

			err = os.RemoveAll("temp/file" + TempFileList[i].FileID)
			ErrCheck(err)
		}
	}
}