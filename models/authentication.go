package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	Id              primitive.ObjectID `json:"id,omitempty"`
	Detail          AuthDetail         `json:"authdetail,omitempty"`
	Hostname        string             `json:"hostname,omitempty"`
	Password        string             `json:"password,omitempty"`
	ContactNumber   string             `json:"contact,omitempty"`
	GoogleAuthKey   string             `json:"googleauth,omitempty"`
	FacebookAuthkey string             `json:"facebookauth,omitempty"`
}

type AuthDetail struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	User        User               `json:"user,omitempty"`
	Status      string             `json:"status,omitempty"`
	AutoLogin   bool               `json:"autologin,omitempty"`
	LastLoginAt time.Time          `json:"lastloginat,omitempty"`
	ApiKey      string             `json:"apikey,omitempty"`
	ApiSecret   string             `json:"apisecret,omitempty"`
}

type Token struct {
	Token string `json:"apiKey,omitempty"`
}
