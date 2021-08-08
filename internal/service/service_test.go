package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/NVTer/rest-api-example/internal"
	errs "github.com/NVTer/rest-api-example/internal/errors"
	"github.com/NVTer/rest-api-example/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	data        *repository.Database   //nolint: gochecknoglobals
	repos       *repository.Repository //nolint: gochecknoglobals
	serv        *Serv                  //nolint: gochecknoglobals
	positionIDs []string               //nolint: gochecknoglobals
	employeeIDs []string               //nolint: gochecknoglobals
)

func initData() {
	data = repository.NewDataBase()
	repos = repository.NewRepo(data)
	serv = NewServ(repos)
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

func createRightContext() context.Context {
	//revive:disable
	return context.WithValue(context.Background(), "correlation_id", uuid.New().String()) //nolint:staticcheck
	//revive:enable
}

func createBadContext() context.Context {
	//revive:disable
	return context.WithValue(context.Background(), "noname_id", uuid.New().String()) //nolint:staticcheck
	//revive:enable
}

func TestCreatePosition(t *testing.T) {
	initData()
	p := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&p)
	newPos := internal.Position{ID: createPosID(), Name: "lead", Salary: decimal.New(2000, 0)}
	testTable := []struct {
		addID          string
		expectedName   string
		expectedSalary decimal.Decimal
		add            internal.Position
		ctx            context.Context
		err            error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			add:            newPos,
			addID:          positionIDs[1],
			expectedName:   newPos.Name,
			expectedSalary: newPos.Salary,
			ctx:            createRightContext(),
			err:            nil,
		},
		{
			add: newPos,
			ctx: createRightContext(),
			err: errs.PositionIsExists(),
		},
	}
	for _, testCase := range testTable {
		id, result := serv.CreatePosition(testCase.ctx, &testCase.add)
		position := repos.GetPositions()[id]
		if result != nil {
			assert.Equal(t, result, testCase.err)
		} else {
			assert.Equal(t, position.Name, testCase.expectedName)
			assert.Equal(t, position.Salary, testCase.expectedSalary)
		}
	}
}

func TestCreateEmployee(t *testing.T) { //nolint:funlen
	initData()
	p := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&p)
	firstEmp := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Vik",
		LasName:    "Sick",
		PositionID: p.ID,
	}
	repos.AddEmployee(&firstEmp)
	newEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Bread",
		LasName:    "Brown",
		PositionID: p.ID,
	}
	fakeEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Big",
		LasName:    "Jo",
		PositionID: uuid.New(),
	}
	testTable := []struct {
		addID             string
		expectedFirstName string
		expectedLasName   string
		positionID        string
		add               internal.Employee
		ctx               context.Context
		err               error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			add:               newEmployee,
			addID:             employeeIDs[1],
			positionID:        positionIDs[0],
			expectedFirstName: newEmployee.FirstName,
			expectedLasName:   newEmployee.LasName,
			ctx:               createRightContext(),
			err:               nil,
		},
		{
			add:               newEmployee,
			addID:             employeeIDs[1],
			positionID:        positionIDs[0],
			expectedFirstName: newEmployee.FirstName,
			expectedLasName:   newEmployee.LasName,
			ctx:               createRightContext(),
			err:               errs.EmployeeIsExists(),
		},
		{
			add:               fakeEmployee,
			addID:             employeeIDs[2],
			expectedFirstName: fakeEmployee.FirstName,
			expectedLasName:   fakeEmployee.LasName,
			positionID:        positionIDs[0],
			ctx:               createRightContext(),
			err:               errs.PositionIsNotExists(),
		},
	}
	for _, testCase := range testTable {
		id, result := serv.CreateEmployee(testCase.ctx, &testCase.add)
		employee := repos.GetEmployees()[id]
		if result != nil {
			assert.Equal(t, result, testCase.err)
		} else {
			assert.Equal(t, employee.FirstName, testCase.expectedFirstName)
			assert.Equal(t, employee.LasName, testCase.expectedLasName)
			assert.Equal(t, employee.PositionID.String(), testCase.positionID)
		}
	}
}

func TestGetPositionZero(t *testing.T) {
	initData()
	testTable := []struct {
		expected []internal.Position
		limit    int
		offset   int
		ctx      context.Context
		err      error
	}{
		{
			expected: []internal.Position{},
			limit:    1,
			offset:   1,
			ctx:      createRightContext(),
			err:      nil,
		},
	}
	for _, testCase := range testTable {
		positions, err := serv.GetPositions(testCase.ctx, testCase.limit, testCase.offset)
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, positions, testCase.expected)
		}
	}
}

func TestGetPositions(t *testing.T) {
	initData()
	firstPosition := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPosition)
	secondPosition := internal.Position{ID: createPosID(), Name: "lead", Salary: decimal.New(2000, 0)}
	repos.AddPosition(&secondPosition)
	testTable := []struct {
		expected []internal.Position
		limit    int
		offset   int
		ctx      context.Context
		err      error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			expected: []internal.Position{
				firstPosition,
				secondPosition,
			},
			limit:  2,
			offset: 1,
			ctx:    createRightContext(),
			err:    nil,
		},
		{
			limit:  120,
			offset: 1,
			ctx:    createRightContext(),
			err:    errs.BadRequest(),
		},
		{
			limit:  1,
			offset: 13,
			ctx:    createRightContext(),
			err:    errs.NotFound(),
		},
	}
	for _, testCase := range testTable {
		positions, err := serv.GetPositions(testCase.ctx, testCase.limit, testCase.offset)
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, positions, testCase.expected)
		}
	}
}

func TestGetEmployeeZero(t *testing.T) {
	initData()
	testTable := []struct {
		expected []internal.Employee
		limit    int
		offset   int
		ctx      context.Context
		err      error
	}{
		{
			expected: []internal.Employee{},
			limit:    1,
			offset:   1,
			ctx:      createRightContext(),
			err:      nil,
		},
	}
	for _, testCase := range testTable {
		employees, err := serv.GetEmployees(testCase.ctx, testCase.limit, testCase.offset)
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, employees, testCase.expected)
		}
	}
}

func TestGetEmployees(t *testing.T) {
	initData()
	firstPosition := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPosition)
	firstEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Nick",
		LasName:    "Bobs",
		PositionID: firstPosition.ID,
	}
	repos.AddEmployee(&firstEmployee)
	secondEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Bob",
		LasName:    "Daddy",
		PositionID: firstPosition.ID,
	}
	repos.AddEmployee(&secondEmployee)
	testTable := []struct {
		expected []internal.Employee
		limit    int
		offset   int
		ctx      context.Context
		err      error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			expected: []internal.Employee{
				firstEmployee,
				secondEmployee,
			},
			limit:  2,
			offset: 1,
			ctx:    createRightContext(),
			err:    nil,
		},
		{
			limit:  120,
			offset: 1,
			ctx:    createRightContext(),
			err:    errs.BadRequest(),
		},
		{
			limit:  1,
			offset: 13,
			ctx:    createRightContext(),
			err:    errs.NotFound(),
		},
	}
	for _, testCase := range testTable {
		employees, err := serv.GetEmployees(testCase.ctx, testCase.limit, testCase.offset)
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, employees, testCase.expected)
		}
	}
}

func TestGetPosition(t *testing.T) {
	initData()
	firstPosition := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPosition)
	secondPosition := internal.Position{ID: createPosID(), Name: "lead", Salary: decimal.New(1500, 0)}
	repos.AddPosition(&secondPosition)

	testTable := []struct {
		expected internal.Position
		id       string
		ctx      context.Context
		err      error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			expected: firstPosition,
			id:       positionIDs[0],
			ctx:      createRightContext(),
			err:      nil,
		},
		{
			ctx: createRightContext(),
			id:  uuid.New().String(),
			err: errs.NotFound(),
		},
		{
			ctx: createRightContext(),
			id:  "9",
			err: errs.BadRequest(),
		},
	}
	for _, testCase := range testTable {
		position, err := serv.GetPosition(testCase.ctx, testCase.id)
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, position, testCase.expected)
		}
	}
}

func TestGetEmployee(t *testing.T) {
	initData()
	firstPosition := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPosition)
	firstEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Nick",
		LasName:    "Bobs",
		PositionID: firstPosition.ID,
	}
	repos.AddEmployee(&firstEmployee)
	testTable := []struct {
		expected internal.Employee
		id       string
		ctx      context.Context
		err      error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			expected: firstEmployee,
			id:       employeeIDs[0],
			ctx:      createRightContext(),
			err:      nil,
		},
		{
			ctx: createRightContext(),
			id:  uuid.New().String(),
			err: errs.NotFound(),
		},
		{
			ctx: createRightContext(),
			id:  "9",
			err: errs.ParseError(),
		},
	}
	for _, testCase := range testTable {
		employee, err := serv.GetEmployee(testCase.ctx, testCase.id)
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, employee, testCase.expected)
		}
	}
}

func TestDeletePosition(t *testing.T) {
	initData()
	firstPosition := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPosition)

	testTable := []struct {
		expected map[string]internal.Position
		delete   string
		ctx      context.Context
		err      error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			expected: map[string]internal.Position{},
			delete:   positionIDs[0],
			ctx:      createRightContext(),
			err:      nil,
		},
	}
	for _, testCase := range testTable {
		err := serv.DeletePosition(testCase.ctx, testCase.delete)
		positions := repos.GetPositions()
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, positions, testCase.expected)
		}
	}
}

func TestDeleteEmployee(t *testing.T) {
	initData()
	firstPosition := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPosition)
	firstEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Nick",
		LasName:    "Bobs",
		PositionID: firstPosition.ID,
	}
	repos.AddEmployee(&firstEmployee)
	testTable := []struct {
		expected map[string]internal.Employee
		delete   string
		ctx      context.Context
		err      error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			expected: map[string]internal.Employee{},
			delete:   employeeIDs[0],
			ctx:      createRightContext(),
			err:      nil,
		},
	}
	for _, testCase := range testTable {
		err := serv.DeleteEmployee(testCase.ctx, testCase.delete)
		employees := repos.GetEmployees()
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else {
			assert.Equal(t, employees, testCase.expected)
		}
	}
}

func TestUpdatePosition(t *testing.T) {
	initData()
	p := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&p)
	updatePos := internal.Position{ID: p.ID, Name: "lead", Salary: decimal.New(2000, 0)}
	posNilID := internal.Position{ID: uuid.Nil, Name: "lead", Salary: decimal.New(2000, 0)}
	testTable := []struct {
		update internal.Position
		ctx    context.Context
		err    error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			update: updatePos,
			ctx:    createRightContext(),
			err:    nil,
		},
		{
			update: posNilID,
			ctx:    createRightContext(),
			err:    errs.BadRequest(),
		},
	}
	for _, testCase := range testTable {
		err := serv.UpdatePosition(testCase.ctx, &testCase.update)
		_, ok := repos.GetPositions()[testCase.update.ID.String()]
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else if !ok {
			t.Errorf("Error!\n Expected : true;\nResult: %#v\n", ok)
		}
	}
}

func TestUpdateEmployee(t *testing.T) {
	initData()
	firstPosition := internal.Position{ID: createPosID(), Name: "worker", Salary: decimal.New(500, 0)}
	repos.AddPosition(&firstPosition)
	firstEmployee := internal.Employee{
		ID:         createEmpID(),
		FirstName:  "Nick",
		LasName:    "Bobs",
		PositionID: firstPosition.ID,
	}
	repos.AddEmployee(&firstEmployee)
	updateEmp := internal.Employee{
		ID:         firstEmployee.ID,
		FirstName:  "Bob",
		LasName:    "Daddy",
		PositionID: firstPosition.ID,
	}
	empNilID := internal.Employee{
		ID:         uuid.Nil,
		FirstName:  "Bb",
		LasName:    "Ddy",
		PositionID: firstPosition.ID,
	}
	testTable := []struct {
		update internal.Employee
		ctx    context.Context
		err    error
	}{
		{
			ctx: createBadContext(),
			err: errs.LogError(),
		},
		{
			update: updateEmp,
			ctx:    createRightContext(),
			err:    nil,
		},
		{
			update: empNilID,
			ctx:    createRightContext(),
			err:    errs.BadRequest(),
		},
	}
	for _, testCase := range testTable {
		err := serv.UpdateEmployee(testCase.ctx, &testCase.update)
		_, ok := repos.GetEmployees()[testCase.update.ID.String()]
		if err != nil {
			assert.Equal(t, err, testCase.err)
		} else if !ok {
			t.Errorf("Error!\n Expected : true;\nResult: %#v\n", ok)
		}
	}
}
