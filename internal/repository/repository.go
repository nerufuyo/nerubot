package repository

import (
	"github.com/nerufuyo/nerubot/internal/pkg/mongodb"
)

// MongoDB is the shared MongoDB client used by all repositories.
// It is set during application startup via SetMongo.
var MongoDB *mongodb.Client

// SetMongo sets the shared MongoDB client for all repositories.
func SetMongo(client *mongodb.Client) {
	MongoDB = client
}
