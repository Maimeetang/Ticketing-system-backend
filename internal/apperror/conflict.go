package apperror

type Conflict interface {
	error
	Conflict() bool
}

type ErrConflict struct {
	Field string
}

func (e *ErrConflict) Error() string {
	return e.Field + " already exists"
}

func (e *ErrConflict) Conflict() bool {
	return true
}

func NewConflict(field string) *ErrConflict {
	return &ErrConflict{Field: field}
}
