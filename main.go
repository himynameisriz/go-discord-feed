package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	// "encoding/json"

	"./RssFeed"
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

	ConfigFile := os.Args[1]
	fmt.Println("Config file found:", ConfigFile)
	filePath, err := filepath.Abs(ConfigFile)
	if err != nil {
		log.Error("Setting file path, ", err)
		return
	}
	fmt.Println("No errors opening file")
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Error("Opening config, ", err)
	} else {
		fmt.Println(jsonFile)
	}
	// json.Unmarshal()
	// Token = os.Args[1]
	// ChannelId = os.Args[2]
	// RssFeed = os.Args[3]

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	botId = u.ID

	err = dg.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
	fmt.Println("Starting RSS Feed")

	// go runRssFeed(dg)
	// <-make(chan struct{})
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
