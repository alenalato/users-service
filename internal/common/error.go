package common

// ErrorType is the type of the error
type ErrorType int

const (
	ErrTypeUnknown ErrorType = iota
	ErrTypeNotFound
	ErrTypeAlreadyExists
	ErrTypeInvalidArgument
	ErrTypeInternal
)

func (e ErrorType) String() string {
	switch e {
	case ErrTypeUnknown:
		return "unknown error"
	case ErrTypeNotFound:
		return "not found"
	case ErrTypeAlreadyExists:
		return "already exists"
	case ErrTypeInvalidArgument:
		return "invalid argument"
	case ErrTypeInternal:
		return "internal error"
	default:
		return "unknown error type"
	}
}

// Error is an auxiliary error type
type Error struct {
	errType ErrorType
	err     error
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) GetType() ErrorType {
	return e.errType
}

func NewError(e error, t ErrorType) Error {
	return Error{err: e, errType: t}
}
