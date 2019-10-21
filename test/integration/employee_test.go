package integration

import (
	"context"
	"fmt"
	"testing"
)

func TestEmployee(t *testing.T) {
	if employeeFilter == nil {
		t.Skip("Skipping test for no employee filter was set.")
	}

	emps, err := tx99.Employee.Find(context.Background(), employeeFilter)
	if err != nil {
		t.Fatalf("got error while calling Employee.Find(%+v): '%s'; want nil.", employeeFilter, err.Error())
	}

	if len(emps) == 0 {
		t.Skip("No Employee found.")
	}

	fmt.Println(emps)
	// employee := r.Employees[0]
}
