package taxis99

import (
	"context"
	"fmt"
	"net/http"
)

const (
	employeesEndpoint           endpoint = `employees`
	employeeEndpoint            endpoint = `employees/%d`
	employeeExternalIdEndpoint  endpoint = `employees/external-id/%d`
	employeeCostCentersEndpoint endpoint = `employees/%d/costcenter`
)

// Hashset for allwed query params.
var employeeFields = map[string]struct{}{
	"search":     struct{}{},
	"limit":      struct{}{},
	"page":       struct{}{},
	"nationalId": struct{}{},
}

type Phone struct {
	Number  string `json:"number,omitempty"`
	Country string `json:"country,omitempty"`
}

type Employee struct {
	ID           int64    `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Email        string   `json:"email,omitempty"`
	Phone        *Phone   `json:"phone,omitempty"`
	Company      *Company `json:"company,omitempty"`
	NationalID   string   `json:"nationalId,omitempty"`
	SupervisorID int64    `json:"supervisorId,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
	ExternalID   int64    `json:"externalId"`
	Categories   []string `json:"categories,omitempty"`
}

type reqEmployee struct {
	Employee         *Employee `json:"employee"`
	SendWelcomeEmail bool      `json:"sendWelcomeEmail"`
}

type EmployeeService service

func (e *EmployeeService) Find(ctx context.Context, f Filter) ([]*Employee, error) {
	var employees []*Employee

	v := f.values(ccFields)

	err := e.client.Request(ctx, http.MethodGet, string(employeesEndpoint.Query(v)), nil, &employees)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (e *EmployeeService) FindByExternalID(ctx context.Context, extID int) ([]*Employee, error) {
	var employees []*Employee

	endpoint := fmt.Sprintf(string(employeeExternalIdEndpoint), extID)

	err := e.client.Request(ctx, http.MethodGet, endpoint, nil, &employees)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (e *EmployeeService) Create(ctx context.Context, emp Employee, sendEmail bool) (*Employee, error) {
	res := new(Employee)

	newEmp := reqEmployee{
		Employee:         &emp,
		SendWelcomeEmail: sendEmail,
	}

	err := e.client.Request(context.Background(), http.MethodPost, string(employeesEndpoint), newEmp, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *EmployeeService) Update(ctx context.Context, emp Employee) (*Employee, error) {
	res := new(Employee)

	updatedEmp := reqEmployee{
		Employee: &emp,
	}

	endpoint := fmt.Sprintf(string(employeeEndpoint), emp.ID)

	err := e.client.Request(context.Background(), http.MethodPut, endpoint, updatedEmp, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (e *EmployeeService) Remove(ctx context.Context, id int64) error {

	endpoint := fmt.Sprintf(string(employeeEndpoint), id)

	return e.client.Request(context.Background(), http.MethodDelete, endpoint, nil, nil)
}

func (e *EmployeeService) FindCostCenters(ctx context.Context, empID int64) ([]*CostCenter, error) {
	var costCenters []*CostCenter
	endpoint := fmt.Sprintf(string(employeeCostCentersEndpoint), empID)

	err := e.client.Request(context.Background(), http.MethodGet, endpoint, nil, &costCenters)
	if err != nil {
		return nil, err
	}

	return costCenters, nil
}

func (e *EmployeeService) UpdateCostCenters(ctx context.Context, empID int64, costCenterIDs []int64) ([]int64, error) {
	type updateCostCenters struct {
		CostCenterIDs []int64 `json:"costCenterIDs"`
	}
	var ids []int64

	endpoint := fmt.Sprintf(string(employeeCostCentersEndpoint), empID)

	newCostCenters := updateCostCenters{costCenterIDs}

	err := e.client.Request(context.Background(), http.MethodPatch, endpoint, newCostCenters, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
