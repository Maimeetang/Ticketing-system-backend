package errors

// ---------- Interfaces ----------

type NotFound interface {
	error
	NotFound() bool
}

type Conflict interface {
	error
	Conflict() bool
}

type BadRequest interface {
	error
	BadRequest() bool
}

// ---------- NotFound ----------

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

// ---------- Conflict ----------

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

// ---------- BadRequest ----------

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