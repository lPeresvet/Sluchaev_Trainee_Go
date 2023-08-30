package errors

import "fmt"

type NotFoundError struct {
	Type string
	Id   string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v with id <%v> not found", e.Type, e.Id)
}

type NotFoundErrorWithMessage struct {
	Message string
}

func (e *NotFoundErrorWithMessage) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
