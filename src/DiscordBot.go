package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var MetaDataChannel string = "874716702988980244"// Replace ID with your channel
var StorageChannel string = "874477682501496852"// Replace ID with your channel

var discord *discordgo.Session

func InitiateBot() {
	// Get discord bot key
	file, err := ioutil.ReadFile("key.txt")// Open the file and insert your own key there
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
		file, err := os.Open("temp/file" + data.FileID + "/chunk" + strconv.Itoa(i))//var reader io.Reader = file
		ErrCheck(err)
		defer file.Close()

		// Send it
		SentFile, err := discord.ChannelFileSend(StorageChannel, "chunk" + strconv.Itoa(i), file)
		ErrCheck(err)

		// Add ID to slice
		AllFiles = append(AllFiles, SentFile.ID)
	}

	// Send meta data
	var ids string
	for i := 0; i < len(AllFiles); i++ {
		ids = ids + "," + AllFiles[i]
	}
	_, err := discord.ChannelMessageSend(MetaDataChannel, data.FileID + "," + data.FileName + ids)
	ErrCheck(err)
}

func DownloadFiles(id string) string {
	// Get last 100 messages
	msgs, err := discord.ChannelMessages(MetaDataChannel, 100, "", "", "")
	ErrCheck(err)

	for i := 0; i < len(msgs); i++ {
		// Check for limit hit (100 msgs)
		if i == 99 {
			msgs2, err := discord.ChannelMessages(MetaDataChannel, 100, msgs[i].ID, "", "")
			ErrCheck(err)

			msgs = msgs2
			i = 0
		}

		// Check ID
		data := strings.Split(msgs[i].Content, ",")
		if data[0] == id {
			// Create directory just in case
			err := os.MkdirAll("public/files", os.ModePerm)
			ErrCheck(err)

			out, err := os.Create("public/files/u" + data[0] + data[1])
			ErrCheck(err)
			defer out.Close()

			// Download all chunks
			for j := 2; j < len(data); j++ {
				fmt.Println("test" +  strconv.Itoa(j))
				msg, err := discord.ChannelMessage(StorageChannel, data[j])
				ErrCheck(err)

				res, err := http.Get(msg.Attachments[0].URL)
				ErrCheck(err)
				defer res.Body.Close()

				_, err = io.Copy(out, res.Body)
				ErrCheck(err)
			}

			// Saved for later cleanup
			TempFileList = append(TempFileList, TempFile{
				data[1],
				data[0],
			})

			return "/files/u" + data[0] + data[1]
		}
	}

	// Couldn't find the file
	return "null"
}
