package ctxRestaurant

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Name          string             `json:"name,omitempty" validate:"requried"`
	Type          string             `json:"type,omitempty" validate:"requried"`
	ContactNumber string             `json:"contact,omitempty"`
	ServiceOption ServiceOption      `json:"serviceoption,omitempty"`
	OpenHours     string             `json:"hours,omitempty"`
	Website       string             `json:"website,omitempty"`
	Address       string             `json:"address,omitempty"`
	Rating        Rating             `json:"rating,omitempty"`
	Menu          Menu               `json:"menu,omitempty"`
}

type Menu struct {
	Id   primitive.ObjectID `json:"id,omitempty"`
	Name string             `json:"name,omitempty" validate:"required"`
	Code string             `json:"code" validate:"required"`
	//pictures
}

type Rating struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	UserId    uint8              `json:"userId,omitempty" validate:"required"`
	CreatedAt time.Time          `json:"createdat,omitempty"`
	UpdatedAt time.Time          `json:"updatedat,omitempty"`
	Status    string             `json:"status,omitempty" validate:"required"`
	Rating    float32            `json:"rating,omitempty" validate:"required"`
	Comment   string             `json:"comment,omitempty"`
}

type ServiceOption struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	OptionName string             `json:"option,omitempty"`
}
