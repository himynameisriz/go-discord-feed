package main

import (
	"encoding/json"
	"fmt"
	io "io/ioutil"
	"os"
	"path/filepath"
	"time"

	// "encoding/json"

	"./rssFeed"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	Token     string
	LastTitle string
	Count     int
	ChannelId string
	RssFeed   string
	LogFile   *os.File
)

var botId string

func init() {
	logPath, _ := filepath.Abs("logs/logs.txt")
	LogFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(LogFile)
}

func main() {
	log.Info("Feed reader started")
	filePath, err := filepath.Abs("config.json")
	if err != nil {
		log.Error("Setting file path, ", err)
		return
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Error("Opening config, ", err)
	}

	defer jsonFile.Close()
	log.Info("No errors opening file")

	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	unmarshalErr := json.Unmarshal(byteValue, &config)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
	}
	log.Info("We have some data", config)

	// // Create a new Discord session using the provided bot token.
	// dg, err := discordgo.New("Bot " + config.Token)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// u, err := dg.User("@me")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// botId = u.ID

	// err = dg.Open()

	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// fmt.Println("Bot is running!")
	// fmt.Println("Starting RSS Feed")

	// // go runRssFeed(dg)
	// // <-make(chan struct{})
	defer LogFile.Sync()
	defer LogFile.Close()
	return
	// Cleanly close down the Discord session.
}

func runRssFeed(s *discordgo.Session) {
	fmt.Println("RSS Feed reader starting")
	for ok := true; ok; ok = (Count < 5) {
		message, lastTitle, err := rssFeed.RunFeed(RssFeed, LastTitle)

		if err != nil {
			Count++
			fmt.Println("Error, ", err)
		} else {
			if len(message) == 0 {
				fmt.Println("No message found, sleep time")
			} else {
				LastTitle = lastTitle
				fmt.Println(message)
				fmt.Println("Attempting to send message")
				s.ChannelMessageSend(ChannelId, message)
				fmt.Println("Message sent")
			}
		}

		fmt.Println("Sleep has began")
		time.Sleep(10 * time.Minute)
		fmt.Println("Sleep ended")
	}
}
