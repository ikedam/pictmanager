package rfc7807

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Error struct {
	Status    int    `json:"status"`
	Type      string `json:"type,omitempty"`
	Title     string `json:"title,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Instance  string `json:"instance"`
	RootCause error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%v: %v", e.Status, e.Detail)
	}
	if e.Title != "" {
		return fmt.Sprintf("%v: %v", e.Status, e.Title)
	}
	if e.Type != "" {
		return fmt.Sprintf("%v: %v", e.Status, e.Type)
	}
	return fmt.Sprintf("status: %v", e.Status)
}

func (e *Error) Unwrap() error {
	return e.RootCause
}

func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') && e.RootCause != nil {
			io.WriteString(s, e.Error())
			fmt.Fprintf(s, "\n%+v", e.RootCause)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func New(status int) *Error {
	return &Error{
		Status: status,
	}
}

func (e *Error) WithType(t string) *Error {
	e.Type = t
	return e
}

func (e *Error) WithTitle(title string) *Error {
	e.Title = title
	return e
}

func (e *Error) WithDetail(detail string) *Error {
	e.Detail = detail
	return e
}

func (e *Error) WithDetailf(format string, args ...any) *Error {
	e.Detail = fmt.Sprintf(format, args...)
	return e
}

func (e *Error) WithInstance(instance string) *Error {
	e.Instance = instance
	return e
}

func (e *Error) WithError(err error) *Error {
	e.RootCause = err
	if e.Detail == "" {
		e.Detail = err.Error()
	}
	return e
}

func As(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func BadRequest() *Error {
	return New(http.StatusBadRequest).WithTitle("Bad Request")
}
