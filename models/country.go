package models

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
	OfficialName string `json:"official_name"`
	Flag string `json:"flag"`
	Capital string `json:"capital"`
	Population string `json:"population"`
	Region string `json:"region"`
	Currency string `json:"currency"`
	Languages string `json:"languages"`
}
