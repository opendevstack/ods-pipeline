package httpjson

import (
	"fmt"
	"net/http"
)

// ProblemDetailer should be implemented by any error that has should
// be encoded as a specific RFC 7807 problem detail.
type ProblemDetailer interface {
	ProblemDetail() any
	Code() int
}

// StatusProblem is basically an error with an HTTP status code.
type StatusProblem struct {
	Status int    `json:"-"`
	Detail string `json:"detail,omitempty"`
	Err    error
}

// NewInternalProblem returns a StatusProblem with code 500.
func NewInternalProblem(detail string, err error) *StatusProblem {
	return NewStatusProblem(http.StatusInternalServerError, detail, err)
}

// NewStatusProblem is a helper to create a StatusProblem.
func NewStatusProblem(code int, detail string, err error) *StatusProblem {
	return &StatusProblem{Status: code, Detail: detail, Err: err}
}

// ProblemDetail turns p into a RFC 7807 problem detail.
func (p *StatusProblem) ProblemDetail() any {
	return &defaultProblem{
		Title:  p.title(),
		Detail: p.Detail,
		Status: p.Status,
		err:    p.Err,
	}
}

// Code returns the HTTP status code of p.
func (p *StatusProblem) Code() int {
	return p.Status
}

func (p *StatusProblem) String() string {
	return fmt.Sprintf("%s: %s: %s", p.title(), p.Detail, p.Err.Error())
}

func (p *StatusProblem) Error() string {
	return p.String()
}

func (p *StatusProblem) title() string {
	return http.StatusText(p.Status)
}

// defaultProblem is a basic RFC 7807 problem detail.
// It wraps an error as that is the typical source of a problem.
type defaultProblem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int    `json:"-"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	err      error
}
