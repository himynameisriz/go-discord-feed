package main

// Config object
type Config struct {
	Token string
	Feeds []Feed
}

//Feed objects
type Feed struct {
	FeedName  string
	ChannelId string
	FeedURL   string
}
