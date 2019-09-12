package taxis99

import (

)

type ExistingUser struct {
	Categories []string `json:"categories"`
	Company    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"company"`
	Email      string `json:"email"`
	Enabled    bool   `json:"enabled"`
	ExternalID int    `json:"externalId"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	NationalID string `json:"nationalId"`
	Phone      struct {
		Country string `json:"country"`
		Number  string `json:"number"`
	} `json:"phone"`
	SupervisorID int `json:"supervisorId"`
}