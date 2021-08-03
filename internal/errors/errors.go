package errors

var (
	badRequest          = newError("bad request")            // nolint: gochecknoglobals
	notFound            = newError("not found")              // nolint: gochecknoglobals
	positionIsExists    = newError("position is exists")     // nolint: gochecknoglobals
	employeeIsExists    = newError("employee is exists")     // nolint: gochecknoglobals
	internalServerError = newError("internal server error")  // nolint: gochecknoglobals
	positionIsNotExists = newError("position is not exists") // nolint: gochecknoglobals
	logError            = newError("log error")              // nolint: gochecknoglobals
	parseError          = newError("parse error")            // nolint: gochecknoglobals
)

type Errors struct {
	description string
}

func newError(desc string) *Errors {
	return &Errors{description: desc}
}

func (m Errors) Error() string {
	return m.description
}

func BadRequest() error {
	return badRequest
}

func NotFound() error {
	return notFound
}

func PositionIsExists() error {
	return positionIsExists
}

func EmployeeIsExists() error {
	return employeeIsExists
}

func StatusInternalServerError() error {
	return internalServerError
}

func PositionIsNotExists() error {
	return positionIsNotExists
}

func LogError() error {
	return logError
}

func ParseError() error {
	return parseError
}
