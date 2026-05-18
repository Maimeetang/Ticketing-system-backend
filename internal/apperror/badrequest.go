package apperror

type BadRequest interface {
	error
	BadRequest() bool
}

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	return e.Message
}

func (e *ErrBadRequest) BadRequest() bool {
	return true
}

func NewBadRequest(msg string) *ErrBadRequest {
	return &ErrBadRequest{Message: msg}
}