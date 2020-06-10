package models

type StringError struct {
	message string
}

func NewStringError(message string) error{
	return &StringError{message}
}

func (e *StringError) Error() string {
	return e.message
}


