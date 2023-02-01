package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Restaurant struct {
	Id                 primitive.ObjectID `json:"id,omitempty"`
	Name               string             `json:"name,omitempty" validate:"requried"`
	Type               string             `json:"type,omitempty" validate:"requried"`
	ContactNumber      string             `json:"contact,omitempty"`
	ServiceOptionEntry uint8              `json:"serviceoption,omitempty"`
	OpenHours          string             `json:"hours,omitempty"`
	Website            string             `json:"website,omitempty"`
	Address            string             `json:"address,omitempty"`
	CommentEntry       uint8              `json:"comment,omitempty"`
	RatingEntry        uint8              `json:"rating,omitempty"`
	ItemsEntry         uint8              `json:"items,omitempty"`
}
