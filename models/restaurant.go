package models

import (
	"time"
)

type Restaurant struct {
	Id            int             `json:"id,omitempty"`
	Name          string          `json:"name,omitempty"`
	Type          string          `json:"type,omitempty"` //foodpanda,grabfood,googlemap,standalone
	ContactNumber string          `json:"contact,omitempty"`
	ServiceOption []ServiceOption `json:"service_option,omitempty"` //dinein,takeaway
	OpenHours     string          `json:"hours,omitempty"`
	Website       string          `json:"website,omitempty"`
	Address       string          `json:"address,omitempty"`
	Status        string          `json:"status,omitempty"` //active, deleted,hided
	Rating        []Rating        `json:"rating,omitempty"`
	Menu          []MenuHeader    `json:"menu,omitempty"`
}

type Rating struct {
	Id        int       `json:"id,omitempty"`
	UserId    uint8     `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Status    string    `json:"status,omitempty"` //active,deleted,hided
	Rating    float32   `json:"rating,omitempty"`
	Comment   string    `json:"comment,omitempty"`
}

type ServiceOption struct {
	OptionName string `json:"option,omitempty"`
}

type MenuHeader struct {
	Id        int           `json:"id,omitempty"`
	MenuKey   int           `json:"menu_key,omitempty"`
	CreatedAt time.Time     `json:"created_at,omitempty"`
	UpdatedAt time.Time     `json:"updated_at,omitempty"`
	Menu      []MenuDetails `json:"menu_details,omitempty"`
}

type MenuDetails struct {
	Id     int     `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Code   string  `json:"code,omitempty"`
	Images []byte  `json:"images,omitempty"`
	Price  float32 `json:"price,omitempty"`
	Status string  `json:"status,omitempty"`
}
