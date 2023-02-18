package models

import (
	"time"
)

type Restaurant struct {
	Id            int              `bson:"id"`
	CreatedAt     time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" bson:"updated_at"`
	Name          string           `json:"name" bson:"name"`
	Type          string           `json:"type" bson:"type"` //foodpanda,grabfood,googlemap,standalone
	ContactNumber string           `json:"contact" bson:"contact"`
	ServiceOption *[]ServiceOption `json:"service_option" bson:"service_option"` //dinein,takeaway
	OpenHours     string           `json:"hours" bson:"hours"`
	Website       string           `json:"website" bson:"website"`
	Address       string           `json:"address" bson:"address"`
	Status        string           `json:"status" bson:"status"` //active, deleted,hided
	Rating        *[]Rating        `json:"rating" bson:"rating"`
	Menu          *Menu            `json:"menu" bson:"menu"`
}

type Rating struct {
	Id        int       `json:"id" bson:"id"`
	UserId    uint8     `json:"user_id" bson:"user_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Status    string    `json:"status" bson:"status"` //active,deleted,hided
	Rating    float32   `json:"rating" bson:"rating"`
	Comment   string    `json:"comment" bson:"comment"`
}

type ServiceOption struct {
	Option string `json:"option" bson:"option"`
}

type Menu struct {
	Id             int           `json:"id" bson:"id"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
	RestaurantName string        `json:"restaurant_name" bson:"restaurant_name"`
	RestaurantId   int           `json:"restaurant_id" bson:"restaurant_id"`
	Menu           []MenuDetails `json:"menu_details" bson:"menu_details"`
}

type MenuDetails struct {
	Name     string  `json:"name" bson:"name"`
	Code     string  `json:"code" bson:"code"`
	Images   []byte  `json:"images" bson:"images"`
	Price    float32 `json:"price" bson:"price"`
	Currency string  `json:"currency" bson:"currency"`
	Status   string  `json:"status" bson:"status"`
}
