package handler

import (
	"context"

	"github.com/NVTer/rest-api-example/internal"
)

type Service interface {
	CreatePosition(ctx context.Context, p *internal.Position) (string, error)
	CreateEmployee(ctx context.Context, e *internal.Employee) (string, error)
	GetPositions(ctx context.Context, limit, offset int) ([]internal.Position, error)
	GetEmployees(ctx context.Context, limit, offset int) ([]internal.Employee, error)
	GetPosition(ctx context.Context, id string) (internal.Position, error)
	GetEmployee(ctx context.Context, id string) (internal.Employee, error)
	DeletePosition(ctx context.Context, id string) error
	DeleteEmployee(ctx context.Context, id string) error
	UpdatePosition(ctx context.Context, p *internal.Position) error
	UpdateEmployee(ctx context.Context, e *internal.Employee) error
}
