package repository

import (
	"github.com/NVTer/rest-api-example/internal"
	"github.com/NVTer/rest-api-example/internal/errors"
)

type Repository struct {
	data *Database
}

func NewRepo(data *Database) *Repository {
	return &Repository{data: data}
}

func (t Repository) GetPositions() map[string]internal.Position {
	return t.data.GetPosition()
}

func (t Repository) GetEmployees() map[string]internal.Employee {
	return t.data.GetEmployees()
}

func (t Repository) AddPosition(p *internal.Position) {
	t.data.GetPosition()[p.ID.String()] = *p
}

func (t Repository) AddEmployee(e *internal.Employee) {
	t.data.GetEmployees()[e.ID.String()] = *e
}

func (t Repository) DeletePosition(id string) error {
	if _, ok := t.data.positions[id]; ok {
		delete(t.data.positions, id)
		return nil
	}
	return errors.NotFound()
}

func (t Repository) DeleteEmployee(id string) error {
	if _, ok := t.data.employees[id]; ok {
		delete(t.data.employees, id)
		return nil
	}
	return errors.NotFound()
}

func (t Repository) UpdatePosition(p *internal.Position) error {
	if _, ok := t.data.positions[p.ID.String()]; ok {
		t.data.positions[p.ID.String()] = *p
		return nil
	}
	return errors.NotFound()
}

func (t Repository) UpdateEmployee(e *internal.Employee) error {
	if _, ok := t.data.employees[e.ID.String()]; ok {
		if _, ok1 := t.data.positions[e.PositionID.String()]; ok1 {
			t.data.employees[e.ID.String()] = *e
			return nil
		}
		return errors.PositionIsNotExists()
	}
	return errors.NotFound()
}

type Database struct {
	employees map[string]internal.Employee
	positions map[string]internal.Position
}

func NewDataBase() *Database {
	return &Database{
		employees: map[string]internal.Employee{},
		positions: map[string]internal.Position{},
	}
}

func (d Database) GetEmployees() map[string]internal.Employee {
	return d.employees
}

func (d Database) GetPosition() map[string]internal.Position {
	return d.positions
}
