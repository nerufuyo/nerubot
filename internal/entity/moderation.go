package entity

import "time"

// Warning represents a moderation warning issued to a user.
type Warning struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	GuildID   string    `json:"guild_id" bson:"guild_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Username  string    `json:"username" bson:"username"`
	Reason    string    `json:"reason" bson:"reason"`
	IssuedBy  string    `json:"issued_by" bson:"issued_by"`
	IssuedAt  time.Time `json:"issued_at" bson:"issued_at"`
}
