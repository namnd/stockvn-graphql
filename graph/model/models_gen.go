// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Company struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	SectorID string `json:"sectorId"`
}

type CompanySearchParams struct {
	Exchange  *string   `json:"exchange"`
	SectorIds []*string `json:"sectorIds"`
}

type NewTodo struct {
	Text   string `json:"text"`
	UserID string `json:"userId"`
}

type Sector struct {
	ID           string  `json:"id"`
	Label        string  `json:"label"`
	LabelEnglish *string `json:"label_english"`
	Exchange     string  `json:"exchange"`
}

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	User *User  `json:"user"`
}

type Trade struct {
	ClosePrice int `json:"closePrice"`
	Volume     int `json:"volume"`
	Timestamp  int `json:"timestamp"`
}

type TradeSearchParams struct {
	Code string `json:"code"`
	From string `json:"from"`
	To   string `json:"to"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
