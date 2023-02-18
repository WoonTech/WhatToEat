package models

import (
	"time"
)

type User struct {
	Id               int       `json:"id" bson:"id"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	Name             string    `json:"name" bson:"name"`
	Type             string    `json:"type" bson:"type"` //google,standalone,facebook
	Status           string    `json:"status" bson:"status"`
	GroupId          *[]uint8  `json:"group_id" bson:"group_id"`
	OwnRestaurant    *[]uint8  `json:"restaurant" bson:"restaurant"`
	ContactNumber    string    `json:"contact" bson:"contact"`
	CurrentLocation  string    `json:"location" bson:"location"`
	Email            string    `json:"email" bson:"email"`
	ChatLogs         *[]uint8  `json:"logs" bson:"logs"`
	PinnedRestaurant *[]uint8  `json:"pinned_restaurant" bson:"pinned_restaurant"`
}
