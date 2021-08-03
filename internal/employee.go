package internal

import "github.com/google/uuid"

type Employee struct {
	ID         uuid.UUID `json:"ID"`
	FirstName  string    `json:"first_name"`
	LasName    string    `json:"las_name"`
	PositionID uuid.UUID `json:"position_id"`
}
