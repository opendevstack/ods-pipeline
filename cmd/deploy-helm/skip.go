package main

// skipRemainingSteps is a pseudo error used to indicate that remaining
// steps should be skipped.
type skipRemainingSteps struct {
	msg string
}

func (e *skipRemainingSteps) Error() string {
	return e.msg
}
