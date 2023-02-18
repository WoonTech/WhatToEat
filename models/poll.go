package models

import (
	"time"
)

type Poll struct {
	Id             int        `json:"id" bson:"id"`
	CreatedAt      time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" bson:"updated_at"`
	ExpiredAt      time.Time  `json:"expired_at" bson:"expired_at"`
	Detail         PollDetail `json:"poll_details" bson:"poll_details"`
	ParticipantsNo uint8      `json:"participants" bson:"participants"`
}

type PollDetail struct {
	Restaurant   Restaurant `json:"restaurant" bson:"restaurant"`
	Participants []User     `json:"participants" bson:"participants"`
	Results      uint8      `json:"vote_result" bson:"results"`
}
