package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	io "io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "encoding/json"

	"./rssFeed"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	logFile *os.File
)

func init() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", os.ModePerm)
	}

	logPath, _ := filepath.Abs("logs/logs.txt")
	LogFile, _ := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(LogFile)
}

func main() {
	defer logFile.Sync()
	defer logFile.Close()

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
	log.Debug("No errors opening file")

	byteValue, _ := io.ReadAll(jsonFile)
	var config Config
	unmarshalErr := json.Unmarshal(byteValue, &config)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr.Error())
		return
	}

	log.Info("Current config ", config)

	// // Create a new Discord session using the provided bot token.
	dg := createAndOpenDiscord(config.Token)
	if dg == nil {
		return
	}

	log.Info("Bot is running!")

	for _, feed := range config.Feeds {
		log.Info("Starting feed for ", feed.FeedName)
		fmt.Println("Starting feed for ", feed.FeedName)
		go runRssFeed(dg, feed)
	}

	<-make(chan struct{})
	return
	// Cleanly close down the Discord session.
}

func runRssFeed(s *discordgo.Session, feed Feed) {
	for {
		fmt.Println(fmt.Sprintf("Reading feed %s, at %s", feed.FeedName, time.Now().Format("1991 Jun 01")))
		message, currentTitle, err := rssFeed.GetLatest(feed.FeedURL)
		if err != nil {
			log.Error("Error, ", err.Error())
		} else {
			sendMessageString(s, feed, currentTitle, message)
		}

		log.Debug("Sleep begin")
		time.Sleep(time.Duration(feed.SleepTime) * time.Second) // Change this to minute
		log.Debug("Sleep end")
	}
}
func sendMessageString(s *discordgo.Session, feed Feed, currentTitle string, messageString string) {
	if len(messageString) == 0 {
		fmt.Println("No message found, sleep time")
		log.Debug("No message found, sleep time")
		return
	}

	lastTitle := getLastTitle(feed.FeedName)
	if strings.Compare(lastTitle, currentTitle) != 0 {
		_, messageErr := s.ChannelMessageSend(feed.ChannelId, messageString)
		if messageErr != nil {
			fmt.Println("Message error", messageErr.Error())
			log.Error(messageErr.Error())
		} else {
			log.Debug("Message sent")
			fmt.Println("Message sent, ", feed.FeedName)
		}

		appendTitle(feed.FeedName, currentTitle)
	} else {
		log.Debug(fmt.Sprintf("Same title found, did not send message: '%s'", lastTitle))
	}
}
func getLastTitle(feedName string) string {

	filePath := fmt.Sprintf("history/%s.txt", feedName)
	absFilePath, _ := filepath.Abs(filePath)
	historyFile, err := os.Open(absFilePath)
	defer historyFile.Close()

	if err != nil {
		log.Error(fmt.Sprintf("%s not found", filePath))
		checkForDirectory("history", os.ModePerm)
		return ""
	}

	return getLastString(historyFile)
}

func getLastString(file *os.File) string {

	var fileStrings []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileStrings = append(fileStrings, scanner.Text())
	}

	lastTitle := fileStrings[len(fileStrings)-1]
	log.Info("Last title found: ", lastTitle)
	return lastTitle
}

func appendTitle(feedName string, title string) {
	historyPath, _ := filepath.Abs(fmt.Sprintf("history/%s.txt", feedName))
	log.Debug("Opening or creating file: ", historyPath)
	historyFile, _ := os.OpenFile(historyPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

	log.Debug("Writing new title:", title)
	historyFile.WriteString(fmt.Sprintf("\r\n%s", title))
	historyFile.Sync()
	historyFile.Close()
	log.Debug(fmt.Sprintf("%s closed", historyPath))
}

func checkForDirectory(directoryName string, mode os.FileMode) {

	directoryPath, _ := filepath.Abs(directoryName)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		log.Debug(fmt.Sprintf("Creating %s directory", directoryName))
		os.Mkdir(directoryName, mode)
	}

}

func createAndOpenDiscord(token string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	_, err = dg.User("@me")
	if err != nil {
		log.Error(err.Error())
	}

	err = dg.Open()
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return dg
}
