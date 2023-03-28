package errors

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type Error struct {
	err   *errors.Error
	stack *stack
}

func NewCode(code int) *Error {
	return &Error{err: errors.New(code, "", "")}
}
func New(code int, reason, message string) *Error {
	return &Error{err: errors.New(code, reason, message)}
}

func Newf(code int, reason, format string, a ...interface{}) *Error {
	return &Error{err: errors.New(code, reason, fmt.Sprintf(format, a...))}
}

func Code(err error) int {
	return errors.Code(err)
}

func Reason(err error) string {
	return errors.Reason(err)
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	e := errors.FromError(err)
	ret := &Error{err: e}
	return ret
}

func (e *Error) Code() int32 {
	return e.err.Code
}
func (e *Error) Reason() string {
	return e.err.Reason
}
func (e *Error) Message() string {
	return e.err.Message
}
func (e *Error) Metadata() map[string]string {
	return e.err.Metadata
}

func (e *Error) WithStack() *Error {
	e1 := e.Clone()
	e1.stack = callers()
	return e1
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: %v stack = %v", e.err, e.stack)
}

func (e *Error) Unwrap() error { return e.err }

func (e *Error) Is(err error) bool {
	if e == nil || e.err == nil {
		return false
	}
	return e.err.Is(err)
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(httpstatus.ToGRPCCode(int(e.err.Code)), e.err.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.err.Reason,
			Metadata: e.err.Metadata,
		})
	return s
}

func (e *Error) Clone() *Error {
	return &Error{err: errors.Clone(e.err), stack: e.stack}
}

func (e *Error) SetMessage(message string) *Error {
	e.err.Message = message
	return e
}

func (e *Error) NewWithMessage(message string) *Error {
	return e.Clone().SetMessage(message)
}
