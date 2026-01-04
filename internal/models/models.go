package models

import "time"

type Order struct {
	ID          int       `json:"ID"`
	UserID      int       `json:"userID"`
	SaveDate    time.Time `json:"saveDate"`
	OrderIssued bool      `json:"orderIssued"`
	IssuedAt    time.Time `json:"issuedAt"`
}
