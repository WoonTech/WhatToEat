package models

type counter struct {
	Restaurant  int `json:"restaurant,omitempty"`
	Rating      int `json:"rating,omitempty"`
	MenuHeader  int `json:"menu_header,omitempty"`
	MenuDetails int `json:"menu_details,omitempty"`
}
