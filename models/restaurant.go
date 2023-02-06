package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Name          string             `json:"name,omitempty"`
	Type          string             `json:"type,omitempty"`
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
	Name string             `json:"name,omitempty"`
	Code string             `json:"code"`
	//pictures
}

type Rating struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	UserId    uint8              `json:"userId,omitempty"`
	CreatedAt time.Time          `json:"createdat,omitempty"`
	UpdatedAt time.Time          `json:"updatedat,omitempty"`
	Status    string             `json:"status,omitempty"`
	Rating    float32            `json:"rating,omitempty"`
	Comment   string             `json:"comment,omitempty"`
}

type ServiceOption struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	OptionName string             `json:"option,omitempty"`
}
