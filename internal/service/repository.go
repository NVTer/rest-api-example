package service

import "github.com/NVTer/rest-api-example/internal"

type Repository interface {
	GetPositions() map[string]internal.Position
	GetEmployees() map[string]internal.Employee
	AddPosition(p *internal.Position)
	AddEmployee(e *internal.Employee)
	DeletePosition(id string) error
	DeleteEmployee(id string) error
	UpdatePosition(p *internal.Position) error
	UpdateEmployee(e *internal.Employee) error
}
