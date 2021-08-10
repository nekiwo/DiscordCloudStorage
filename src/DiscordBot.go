package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var discord *discordgo.Session

func InitiateBot() {
	// Get discord bot key
	file, err := ioutil.ReadFile("key.txt")
	ErrCheck(err)

	// Discord auth
	discord, err = discordgo.New("Bot " + string(file))
	ErrCheck(err)

	// Idk what this does tbh
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	// Websocket stuff
	err = discord.Open()
	ErrCheck(err)

	// Safe closing
	fmt.Println("Press CTRL + C to close")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
	os.Exit(0)
}

func UploadFiles(data MetaData) {
	// Slice of message IDs
	var AllFiles []string

	// Upload each chunk
	for i := 0; i < data.TotalChunks; i++ {
		// Open file
		file, err := os.Open("./temp/file" + data.FileID + "/chunk" + strconv.Itoa(i))//var reader io.Reader = file
		ErrCheck(err)
		defer file.Close()

		// Send it
		SentFile, err := discord.ChannelFileSend("874477682501496852", "chunk" + strconv.Itoa(i), file)
		ErrCheck(err)

		// Add ID to slice
		AllFiles = append(AllFiles, SentFile.ID)
	}

	// Send meta data
	var ids string
	for i := 0; i < len(AllFiles); i++ {
		ids = ids + "," + AllFiles[i]
	}
	discord.ChannelMessageSend("874716702988980244", data.FileName + "," + data.FileID + ids)
}

func DownloadFiles(id string) {

}