package ctxUser

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id               primitive.ObjectID `json:"id,omitempty"`
	GroupId          []uint8            `json:"groupid,omitempty"`
	Type             string             `json:"type,omitempty"`
	CreatedAt        time.Time          `json:"createdat,omitempty"`
	UpdatedAt        time.Time          `json:"updatedat,omitempty"`
	OwnRestaurant    []uint8            `json:"restaurant,omitempty"`
	ContactNumber    string             `json:"contact,omitempty"`
	Name             string             `json:"name,omitempty"`
	CurrentLocation  string             `json:"location,omitempty"`
	Email            string             `json:"email,omitempty"`
	ChatLogs         []uint8            `json:"logs,omitempty"`
	PinnedRestaurant []uint8            `json:"pinnedres,omitempty"`
}
