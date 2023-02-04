package ctxRestaurant

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
	MenuEntry          Items              `json:"items,omitempty"`
}

type Items struct {
	Id   primitive.ObjectID `json:"id,omitempty"`
	Name string             `json:"name,omitempty" validate:"required"`
	Code string             `json:"code" validate:"required"`
}

type Rating struct {
	Id     primitive.ObjectID `json:"id,omitempty"`
	Rating float32            `json:"rating,omitempty"`
}
