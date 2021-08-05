package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/NVTer/rest-api-example/internal"
	errs "github.com/NVTer/rest-api-example/internal/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	data        *Database   //nolint:gochecknoglobals
	repos       *Repository //nolint:gochecknoglobals
	positionIDs []string    //nolint:gochecknoglobals
	employeeIDs []string    //nolint:gochecknoglobals
)

func updateData() {
	data = NewDataBase()
	repos = NewRepo(data)
	positionIDs = make([]string, 0)
	employeeIDs = make([]string, 0)
}

func createPosID() uuid.UUID {
	id := uuid.New()
	positionIDs = append(positionIDs, id.String())
	return id
}

func createEmpID() uuid.UUID {
	id := uuid.New()
	employeeIDs = append(employeeIDs, id.String())
	return id
}

func TestGetPositions(t *testing.T) {
	updateData()
	p := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&p)
	testTable := []struct {
		expected map[string]internal.Position
	}{
		{
			expected: map[string]internal.Position{
				positionIDs[0]: p,
			},
		},
	}
	for _, testCase := range testTable {
		result := repos.GetPositions()
		assert.Equal(t, result, testCase.expected)
	}
}

func TestGetEmployees(t *testing.T) {
	updateData()
	p := internal.Position{
		ID:     createPosID(),
		Name:   "worker",
		Salary: decimal.New(500, 0),
	}
	repos.AddPosition(&p)
	e := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Nick",
		LasName:    "Bobs",
		PositionID: p.ID,
	}
	repos.AddEmployee(&e)
	testTable := []struct {
		expected map[string]internal.Employee
	}{
		{
			expected: map[string]internal.Employee{
				employeeIDs[0]: e,
			},
		},
	}
	for _, testCase := range testTable {
		result := repos.GetEmployees()
		assert.Equal(t, result, testCase.expected)
	}
}

func TestAddEmployee(t *testing.T) {
	updateData()
	p := internal.Position{
		ID:     createPosID(),
		Name:   "worker",
		Salary: decimal.New(500, 0),
	}
	repos.AddPosition(&p)
	e := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Nick",
		LasName:    "Bobs",
		PositionID: p.ID,
	}
	repos.AddEmployee(&e)
	newEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Bread",
		LasName:    "Brown",
		PositionID: p.ID,
	}
	testTable := []struct {
		expected map[string]internal.Employee
		add      internal.Employee
	}{
		{
			add: e,
			expected: map[string]internal.Employee{
				employeeIDs[0]: e,
			},
		},
		{
			add: newEmployee,
			expected: map[string]internal.Employee{
				employeeIDs[1]: newEmployee,
			},
		},
	}
	for _, testCase := range testTable {
		updateData()
		repos.AddEmployee(&testCase.add)
		result := repos.GetEmployees()
		assert.Equal(t, result, testCase.expected)
	}
}

func TestAddPosition(t *testing.T) {
	updateData()
	p := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&p)
	newPos := internal.Position{ID: createPosID(), Name: "lead", Salary: decimal.New(2000, 0)}
	testTable := []struct {
		expected map[string]internal.Position
		add      internal.Position
	}{
		{
			add: p,
			expected: map[string]internal.Position{
				positionIDs[0]: p,
			},
		},
		{
			add: newPos,
			expected: map[string]internal.Position{
				positionIDs[1]: newPos,
			},
		},
	}
	for _, testCase := range testTable {
		updateData()
		repos.AddPosition(&testCase.add)
		result := repos.GetPositions()
		assert.Equal(t, result, testCase.expected)
	}
}

func TestDeleteEmployee(t *testing.T) {
	updateData()
	p := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&p)
	firstEmp := internal.Employee{
		ID: createEmpID(), FirstName: "Nick",
		LasName:    "Bobs",
		PositionID: p.ID,
	}
	repos.AddEmployee(&firstEmp)
	secondEmp := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Bread",
		LasName:    "Brown",
		PositionID: p.ID,
	}
	repos.AddEmployee(&secondEmp)
	fakeEmp := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Blue",
		LasName:    "Bluer",
		PositionID: p.ID,
	}
	testTable := []struct {
		expected map[string]internal.Employee
		delete   internal.Employee
		err      error
	}{
		{
			delete: firstEmp,
			expected: map[string]internal.Employee{
				employeeIDs[1]: secondEmp,
			},
			err: nil,
		},
		{
			delete:   fakeEmp,
			expected: map[string]internal.Employee{},
			err:      errs.NotFound(),
		},
	}
	for _, testCase := range testTable {
		result := repos.DeleteEmployee(testCase.delete.ID.String())
		if result != nil {
			assert.Equal(t, result, testCase.err)
		} else {
			assert.Equal(t, repos.GetEmployees(), testCase.expected)
		}
	}
}

func TestDeletePosition(t *testing.T) {
	updateData()
	firstPos := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPos)
	secondPos := internal.Position{ID: createPosID(), Name: "lead", Salary: decimal.New(2000, 0)}
	repos.AddPosition(&secondPos)
	fakePos := internal.Position{ID: createPosID(), Name: "principal", Salary: decimal.New(4500, 0)}
	testTable := []struct {
		expected map[string]internal.Position
		delete   internal.Position
		err      error
	}{
		{
			delete: firstPos,
			expected: map[string]internal.Position{
				positionIDs[1]: secondPos,
			},
			err: nil,
		},
		{
			delete:   fakePos,
			expected: map[string]internal.Position{},
			err:      errs.NotFound(),
		},
	}
	for _, testCase := range testTable {
		result := repos.DeletePosition(testCase.delete.ID.String())
		if result != nil {
			assert.Equal(t, result, testCase.err)
		} else {
			assert.Equal(t, repos.GetPositions(), testCase.expected)
		}
	}
}

func TestUpdateEmployee(t *testing.T) {
	updateData()
	p := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&p)
	firstEmp := internal.Employee{ID: createEmpID(), FirstName: "Nick", LasName: "Bobs", PositionID: p.ID}
	repos.AddEmployee(&firstEmp)
	id, err := uuid.Parse(employeeIDs[0])
	if err != nil {
		t.Error(err)
	}
	updateEmp := internal.Employee{ID: id, FirstName: "Bread", LasName: "Brown", PositionID: p.ID}
	repos.AddEmployee(&updateEmp)
	fakeEmp := internal.Employee{ID: createEmpID(), FirstName: "Blue", LasName: "Bluer", PositionID: p.ID}
	noPosEmp := internal.Employee{ID: id, FirstName: "James", LasName: "White", PositionID: uuid.New()}
	testTable := []struct {
		expectedFirstName string
		id                uuid.UUID
		expectedLasName   string
		expectedPosition  internal.Position
		update            internal.Employee
		err               error
	}{
		{
			update:            updateEmp,
			id:                id,
			expectedFirstName: "Bread",
			expectedLasName:   "Brown",
			expectedPosition:  p,
			err:               nil,
		},
		{
			update:            fakeEmp,
			expectedFirstName: "",
			expectedLasName:   "",
			expectedPosition:  internal.Position{},
			err:               errs.NotFound(),
		},
		{
			update:            noPosEmp,
			expectedFirstName: "",
			expectedLasName:   "",
			expectedPosition:  internal.Position{},
			err:               errs.PositionIsNotExists(),
		},
	}
	for _, testCase := range testTable {
		result := repos.UpdateEmployee(&testCase.update)
		if result != nil {
			assert.Equal(t, result, testCase.err)
		} else {
			assert.Equal(t, repos.GetEmployees()[id.String()].FirstName, testCase.expectedFirstName)
			assert.Equal(t, repos.GetEmployees()[id.String()].LasName, testCase.expectedLasName)
			assert.Equal(t, repos.GetEmployees()[id.String()].PositionID, testCase.expectedPosition.ID)
		}
	}
}

func TestUpdatePosition(t *testing.T) {
	updateData()
	firstPos := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPos)
	id, err := uuid.Parse(positionIDs[0])
	if err != nil {
		t.Error(err)
	}
	updatePos := internal.Position{ID: id, Name: "principal", Salary: decimal.New(4500, 0)}
	fakePos := internal.Position{ID: createPosID(), Name: "lead", Salary: decimal.New(2000, 0)}
	testTable := []struct {
		expectedName   string
		id             uuid.UUID
		expectedSalary decimal.Decimal
		update         internal.Position
		err            error
	}{
		{
			update:         updatePos,
			id:             id,
			expectedName:   "principal",
			expectedSalary: decimal.New(4500, 0),
			err:            nil,
		},
		{
			update:         fakePos,
			expectedName:   "",
			expectedSalary: decimal.Zero,
			err:            errs.NotFound(),
		},
	}
	for _, testCase := range testTable {
		result := repos.UpdatePosition(&testCase.update)
		if result != nil {
			assert.Equal(t, result, testCase.err)
		} else {
			assert.Equal(t, repos.GetPositions()[id.String()].Name, testCase.expectedName)
			assert.Equal(t, repos.GetPositions()[id.String()].Salary, testCase.expectedSalary)
		}
	}
}
