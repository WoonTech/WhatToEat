package models

import (
	"time"
)

type Session struct {
	Id        int       `json:"id" bson:"id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Username  string    `json:"username" bson:"username"`
	SessionId string    `json:"session_id" bson:"session_id"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
	Status    string    `json:"status" bson:"status"`
}

type Credentials struct {
	Id            int        `json:"id" bson:"id"`
	CreatedAt     time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" bson:"updated_at"`
	Username      string     `json:"username" bson:"username"`
	Password      string     `json:"password" bson:"password"`
	ContactNumber string     `json:"contact" bson:"contact"`
	Email         string     `json:"email" bson:"email"`
	Detail        CredDetail `json:"auth_details" bson:"auth_details"`
}

type CredDetail struct {
	User        User      `json:"user" bson:"user"`
	Status      string    `json:"status" bson:"status"`
	LastLoginAt time.Time `json:"login_at" bson:"login_at"`
}
