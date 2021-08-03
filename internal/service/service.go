package service

import (
	"context"

	"github.com/NVTer/rest-api-example/internal"
	"github.com/NVTer/rest-api-example/internal/errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Serv struct {
	repo Repository
}

func NewServ(repository Repository) *Serv {
	return &Serv{
		repo: repository,
	}
}

func logCorrelationID(ctx context.Context) error {
	correlationID, ok := ctx.Value("correlation_id").(string)
	if !ok {
		return errors.StatusInternalServerError()
	}
	logrus.WithFields(logrus.Fields{
		"correlation_id": correlationID,
	}).Info()
	return nil
}

func (t Serv) CreatePosition(ctx context.Context, p *internal.Position) (string, error) {
	err := logCorrelationID(ctx)
	if err != nil {
		return "", errors.LogError()
	}
	m := t.repo.GetPositions()
	for _, value := range m {
		if value.Salary.String() == p.Salary.String() && value.Name == p.Name {
			return "", errors.PositionIsExists()
		}
	}
	p.ID = uuid.New()
	t.repo.AddPosition(p)
	return p.ID.String(), nil
}

func (t Serv) CreateEmployee(ctx context.Context, e *internal.Employee) (string, error) {
	err := logCorrelationID(ctx)
	if err != nil {
		return "", errors.LogError()
	}
	m := t.repo.GetEmployees()
	p := t.repo.GetPositions()
	ok := false
	for _, value := range p {
		if value.ID == e.PositionID {
			ok = true
			break
		}
	}
	if !ok {
		return "", errors.PositionIsNotExists()
	}
	for _, value := range m {
		if value.LasName == e.LasName &&
			value.FirstName == e.FirstName {
			return "", errors.EmployeeIsExists()
		}
	}
	e.ID = uuid.New()
	t.repo.AddEmployee(e)
	return e.ID.String(), nil
}

func (t Serv) GetPositions(ctx context.Context, limit, offset int) ([]internal.Position, error) {
	if limit > 100 {
		return nil, errors.BadRequest()
	}
	err := logCorrelationID(ctx)
	if err != nil {
		return nil, errors.LogError()
	}
	m := t.repo.GetPositions()
	answer := make([]internal.Position, 0)
	if len(m) == 0 && offset == 1 && limit == 1 {
		return answer, nil
	}
	positions := make([]internal.Position, 0)
	for _, value := range m {
		positions = append(positions, value)
	}
	offset--
	if float64(len(positions))/float64(limit) <= float64(offset) || limit < 1 || offset < 0 {
		return nil, errors.NotFound()
	}
	for i := limit * offset; i < limit*offset+limit && i < len(positions); i++ {
		answer = append(answer, positions[i])
	}
	return answer, nil
}

func (t Serv) GetEmployees(ctx context.Context, limit, offset int) ([]internal.Employee, error) {
	if limit > 100 {
		return nil, errors.BadRequest()
	}
	err := logCorrelationID(ctx)
	if err != nil {
		return nil, errors.LogError()
	}
	m := t.repo.GetEmployees()
	answer := make([]internal.Employee, 0)
	if len(m) == 0 && offset == 1 && limit == 1 {
		return answer, nil
	}
	employees := make([]internal.Employee, 0)
	for _, value := range m {
		employees = append(employees, value)
	}
	offset--
	if float64(len(employees))/float64(limit) <= float64(offset) || limit < 1 || offset < 0 {
		return nil, errors.NotFound()
	}
	for i := limit * offset; i < limit*offset+limit && i < len(employees); i++ {
		answer = append(answer, employees[i])
	}
	return answer, nil
}

func (t Serv) GetPosition(ctx context.Context, id string) (internal.Position, error) {
	err := logCorrelationID(ctx)
	if err != nil {
		return internal.Position{}, errors.LogError()
	}
	m := t.repo.GetPositions()
	uID, err := uuid.Parse(id)
	if err != nil {
		return internal.Position{}, errors.BadRequest()
	}
	for _, value := range m {
		if value.ID == uID {
			return value, nil
		}
	}
	return internal.Position{}, errors.NotFound()
}

func (t Serv) GetEmployee(ctx context.Context, id string) (internal.Employee, error) {
	err := logCorrelationID(ctx)
	if err != nil {
		return internal.Employee{}, errors.LogError()
	}
	m := t.repo.GetEmployees()
	uID, err := uuid.Parse(id)
	if err != nil {
		return internal.Employee{}, errors.ParseError()
	}
	for _, value := range m {
		if value.ID == uID {
			return value, nil
		}
	}
	return internal.Employee{}, errors.NotFound()
}

func (t Serv) DeletePosition(ctx context.Context, id string) error {
	err := logCorrelationID(ctx)
	if err != nil {
		return errors.LogError()
	}
	return t.repo.DeletePosition(id)
}

func (t Serv) DeleteEmployee(ctx context.Context, id string) error {
	err := logCorrelationID(ctx)
	if err != nil {
		return errors.LogError()
	}
	return t.repo.DeleteEmployee(id)
}

func (t Serv) UpdatePosition(ctx context.Context, p *internal.Position) error {
	err := logCorrelationID(ctx)
	if err != nil {
		return errors.LogError()
	}
	if p.ID.String() == uuid.Nil.String() {
		return errors.BadRequest()
	}
	return t.repo.UpdatePosition(p)
}

func (t Serv) UpdateEmployee(ctx context.Context, e *internal.Employee) error {
	err := logCorrelationID(ctx)
	if err != nil {
		return errors.LogError()
	}
	if e.ID == uuid.Nil {
		return errors.BadRequest()
	}
	return t.repo.UpdateEmployee(e)
}
