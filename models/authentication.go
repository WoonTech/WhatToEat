package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CredDetail struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	UserId      primitive.ObjectID `json:"user,omitempty"`
	Status      string             `json:"status,omitempty"`
	LastLoginAt time.Time          `json:"lastloginat,omitempty"`
}

type Session struct {
	Username  string    `json:"username,omitempty"`
	SessionId string    `json:"sessionId,omitempty"`
	ExpiredAt time.Time `json:"expiredat,omitempty"`
}

type Credentials struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Detail        CredDetail         `json:"authdetail,omitempty"`
	Username      string             `json:"username"`
	Password      string             `json:"password"`
	ContactNumber string             `json:"contact,omitempty"`
}
