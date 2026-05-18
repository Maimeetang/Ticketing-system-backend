package apperror

type NotFound interface {
	error
	NotFound() bool
}

type ErrNotFound struct {
	Object string
}

func (e *ErrNotFound) Error() string {
	return e.Object + " not found"
}

func (e *ErrNotFound) NotFound() bool {
	return true
}

func NewNotFound(object string) *ErrNotFound {
	return &ErrNotFound{Object: object}
}