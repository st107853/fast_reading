package services

import (
	"errors"
	"fmt"
	"net/http"
)

// DomainError carries both a user-safe message and an HTTP status,
// without importing any web framework.
type DomainError struct {
	Code        int
	UserMsg     string
	InternalErr error
}

func (e *DomainError) Error() string {
	if e.InternalErr != nil {
		return fmt.Sprintf("domain error %d: %v", e.Code, e.InternalErr)
	}
	return fmt.Sprintf("domain error %d: %s", e.Code, e.UserMsg)
}

func (e *DomainError) Unwrap() error {
	return e.InternalErr
}

// Constructors for common cases
func ErrNotFound(resource string, internalErr error) *DomainError {
	return &DomainError{
		Code:        http.StatusNotFound,
		UserMsg:     fmt.Sprintf("%s not found.", resource),
		InternalErr: internalErr,
	}
}

func ErrForbidden(action string) *DomainError {
	return &DomainError{
		Code:    http.StatusForbidden,
		UserMsg: fmt.Sprintf("You do not have permission to %s.", action),
	}
}

func ErrInternal(internalErr error) *DomainError {
	return &DomainError{
		Code:        http.StatusInternalServerError,
		UserMsg:     "Something went wrong. Please try again later.",
		InternalErr: internalErr,
	}
}

func ErrUnauthorized() *DomainError {
	return &DomainError{
		Code:    http.StatusUnauthorized,
		UserMsg: "You must be logged in to do that.",
	}
}

func ErrServiceNotInitialized() *DomainError {
	return &DomainError{
		Code:    http.StatusInternalServerError,
		UserMsg: "Service is not available",
	}
}

func ErrBadRequest(userMsg string) *DomainError {
	return &DomainError{
		Code:    http.StatusBadRequest,
		UserMsg: userMsg,
	}
}

func ErrDomainWithMsg(userMsg string, internalErr error) *DomainError {
	return &DomainError{
		Code:        http.StatusInternalServerError,
		UserMsg:     userMsg,
		InternalErr: internalErr,
	}
}

func ErrUser(userMsg string, internalErr error) *DomainError {
	return &DomainError{
		Code:        http.StatusBadRequest,
		UserMsg:     userMsg,
		InternalErr: internalErr,
	}
}

// IsDomainError lets callers check and unwrap in one step
func IsDomainError(err error) (*DomainError, bool) {
	var de *DomainError
	ok := errors.As(err, &de)
	return de, ok
}
