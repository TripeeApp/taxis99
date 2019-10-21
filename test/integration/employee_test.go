package integration

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/mobilitee-smartmob/taxis99"
)

func TestEmployee(t *testing.T) {
	emp := getEmployee(t)
	if emp == nil {
		t.Skip("Employee not found. Skipping test.")
		return
	}
	oldName := emp.Name
	defer func() {
		if employeeFilter != nil {
			emp.Name = oldName
			_, err := tx99.Employee.Update(context.Background(), *emp)
			if err != nil {
				t.Errorf("Got error calling Employee.Update(%+v): '%s'; want nil.", emp, err.Error())
			}
		} else {
			err := tx99.Employee.Remove(context.Background(), emp.ID)
			if err != nil {
				t.Errorf("Got error calling Employee.Remove(%d): '%s'; want nil.", emp.ID, err.Error())
			}
		}
	}()

	emp.Name = randString(10, letterBytes)

	newEmp, err := tx99.Employee.Update(context.Background(), *emp)
	if err != nil {
		t.Fatalf("Got error calling Employee.Update(%+v): '%s'; want nil.", emp, err.Error())
	}

	// Searchs for the new name.
	foundEmp, err := tx99.Employee.Find(context.Background(), taxis99.Filter{"search": newEmp.Name})
	if err != nil {
		t.Fatalf("Got error calling Employee.Find(%+v): '%s'; want nil.", newEmp.Name, err.Error())
	}

	if len(foundEmp) == 0 {
		t.Errorf("Got employee %+v not found; want it to be found.", newEmp)
	}
}

func TestEmployeeCostCenter(t *testing.T) {
	emp := getEmployee(t)
	if emp == nil {
		t.Skip("Employee not found. Skipping test.")
		return
	}
	defer func() {
		if employeeFilter == nil {
			err := tx99.Employee.Remove(context.Background(), emp.ID)
			if err != nil {
				t.Errorf("Got error calling Employee.Remove(%d): '%s'; want nil.", emp.ID, err.Error())
			}
		}
	}()

	newCC := taxis99.CostCenter{Name: randString(10, letterBytes)}

	cc, err := tx99.CostCenter.Create(context.Background(), newCC)
	if err != nil {
		t.Fatalf("Got error calling CostCenter.Create(%+v): '%s'; want nil.", newCC, err.Error())
	}
	defer func() {
		err := tx99.CostCenter.Remove(context.Background(), cc.ID)
		if err != nil {
			t.Errorf("Got error removing cost center %s: %s; want nil.", cc.Name, err.Error())
		}
	}()

	empCCs, err := tx99.Employee.FindCostCenters(context.Background(), emp.ID)
	if err != nil {
		t.Fatalf("Got error calling Employee.FindCostCenters(%d): '%s'; want nil.", emp.ID, err.Error())
	}

	var newEmpCCs []int64

	for _, empCC := range empCCs {
		newEmpCCs = append(newEmpCCs, empCC.ID)
	}

	// Append the new cost center.
	newEmpCCs = append(newEmpCCs, cc.ID)

	_, err = tx99.Employee.UpdateCostCenters(context.Background(), emp.ID, newEmpCCs)
	if err != nil {
		t.Fatalf("Got error calling Employee.UpdateCostCenters(%+v): '%s'; want nil.", newEmpCCs, err.Error())
	}
}

func getEmployee(t *testing.T) *taxis99.Employee {
	if employeeFilter == nil {
		name := randString(10, letterBytes)
		externalID, _ := strconv.Atoi(randString(10, numberBytes))
		newEmp := taxis99.Employee{
			Name:  name,
			Email: fmt.Sprintf("%s@test.com", name),
			Phone: &taxis99.Phone{
				Number:  "11999999999",
				Country: "BRA",
			},
			ExternalID: int64(externalID),
			NationalID: randString(11, numberBytes),
			Enabled:    true,
			Categories: []string{"regular-taxi", "turbo-taxi", "pop99"},
		}

		emp, err := tx99.Employee.Create(context.Background(), newEmp, false)
		if err != nil {
			t.Errorf("Got error creating employee %+v: %s; want nil.", newEmp, err.Error())
			return nil
		}
		return emp
	}

	emps, err := tx99.Employee.Find(context.Background(), employeeFilter)
	if err != nil {
		t.Fatalf("Got error calling Employee.Find(%+v): '%s'; want nil.", employeeFilter, err.Error())
	}

	if len(emps) == 0 {
		t.Errorf("Employee not found.")
		return nil
	}
	return emps[0]
}
