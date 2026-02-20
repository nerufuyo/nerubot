package entity

import "time"

// DadJoke represents a dad joke.
type DadJoke struct {
	ID        string    `json:"id" bson:"id"`
	Setup     string    `json:"setup" bson:"setup"`         // empty for single-line jokes
	Punchline string    `json:"punchline" bson:"punchline"` // the full joke text if single-line
	Source    string    `json:"source" bson:"source"`
	FetchedAt time.Time `json:"fetched_at" bson:"fetched_at"`
}

// Meme represents a meme fetched from the internet.
type Meme struct {
	Title     string    `json:"title" bson:"title"`
	URL       string    `json:"url" bson:"url"`           // direct image URL
	PostLink  string    `json:"post_link" bson:"post_link"` // original post link
	Subreddit string    `json:"subreddit" bson:"subreddit"`
	Author    string    `json:"author" bson:"author"`
	NSFW      bool      `json:"nsfw" bson:"nsfw"`
	FetchedAt time.Time `json:"fetched_at" bson:"fetched_at"`
}
