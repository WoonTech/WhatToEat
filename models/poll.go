package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Poll struct {
	Id             primitive.ObjectID `json:"id,omitempty"`
	Detail         PollDetail         `json:"detail,omitempty"`
	ParticipantsNo uint8              `json:"No,omitempty"`
	CreatedAt      time.Time          `json:"createdat,omitempty"`
	UpdatedAt      time.Time          `json:"updatedat,omitempty"`
	ExpiredAt      time.Time          `json:"expiredat,omitempty"`
}

type PollDetail struct {
	Id           primitive.ObjectID `json:"id,omitempty"`
	Restaurant   Restaurant         `json:"restaurant,omitempty"`
	Participants []User             `json:"participants,omitempty"`
	Results      uint8              `json:"results,omitempty"`
}
