package rssFeed

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

// Get the latest feed item, embedded
func GetLatestEmbed(rssFeed string) (Message, string, error) {
	fp := gofeed.NewParser()
	fmt.Println("Feed created")

	feed, err := fp.ParseURL(rssFeed)
	if err != nil {
		fmt.Println("Error found, ", err)
		return Message{}, "", err
	}
	fmt.Println("Feed read", time.Now())
	rssMessage := createMessage(feed.Items[0].Title, feed.Items[0].Link)
	fmt.Println("Message created, ", rssMessage)
	return rssMessage, feed.Items[0].Title, nil
}

// Get the latest feed item as a string
func GetLatest(rssFeed string) (string, string, error) {
	fp := gofeed.NewParser()
	fmt.Println("Feed created")

	feed, err := fp.ParseURL(rssFeed)
	if err != nil {
		fmt.Println("Error found, ", err)
		return "", "", err
	}
	fmt.Println("Feed read", time.Now())
	messageString := fmt.Sprintf("%s\r\n%s", feed.Items[0].Title, feed.Items[0].Link)
	fmt.Println("Message created, ", messageString)
	return messageString, feed.Items[0].Title, nil
}

// Message object to send embedded messages
type Message struct {
	Title string
	URL   string
}

func createMessage(title string, url string) Message {
	return Message{title, url}
}
