package rssFeed

import (
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// Run Feed to do stuff
func RunFeed(rssFeed string, lastTitle string) (string, string, error) {
	fmt.Println("We are in the feed now")
	fp := gofeed.NewParser()
	fmt.Println("Feed created")

	feed, err := fp.ParseURL(rssFeed)
	fmt.Println("Feed read")
	fmt.Println("Count of items: ", len(feed.Items))

	if err != nil {
		fmt.Println("Error found, ", err)
		return "", lastTitle, err
	}

	if strings.Compare(lastTitle, feed.Items[0].Title) != 0 {
		messageString := feed.Items[0].Title + "\r\n" + feed.Items[0].Link
		return messageString, feed.Items[0].Title, nil
	} else {
		s := fmt.Sprintf("- same title found, not posting '%s'", lastTitle)
		fmt.Println(time.Now(), s)
	}

	return "", lastTitle, nil

}
