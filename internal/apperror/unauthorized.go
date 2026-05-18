package apperror

type Unauthorized interface {
	error
	Unauthorized() bool
}

type ErrUnauthorized struct {
	Message string
}

func (e *ErrUnauthorized) Error() string {
	return e.Message
}

func (e *ErrUnauthorized) Unauthorized() bool {
	return true
}

func NewUnauthorized(msg string) *ErrUnauthorized {
	return &ErrUnauthorized{Message: msg}
}