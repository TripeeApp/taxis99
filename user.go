package taxis99

import (
	"fmt"
)

// User struct to describe an existing user on the 99 taxis API
type User struct {
	Name string `json:"name"`
	Phone struct {
		Number string `json:"number"`
		Country string `json:"country,omitempty"`
	} `json:"phone"`
	Email string `json:"email"`
	NationalID string `json:"nationalId,omitempty"`
	ExternalID string `json:"externalId,omitempty"`
	SupervisorID int `json:"supervisorId,omitempty"`
	Enabled bool `json:"enabled,omitempty"`
	Categories []string `json:"categories"`
	ID int `json:"id,omitempty"`
}

// employee is not exported because it is only present when inserting users
// in order to keep a sane API we use it behind the scenes without exporting it
type employee struct {
	Employee User `json:"employee"`
	SendWelcomeEmail bool `json:"sendWelcomeEmail,omitempty"`
}

// CostCenter describes a cost center in the 99 taxis API
type CostCenter struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

// Company represents a company for marshalling purposes
type Company struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

var (
	listCompanies endpoint = `companies`
	employees endpoint = `employees`
	costCenter = `/costcenter`
)