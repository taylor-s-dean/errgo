package errgo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
)

type ErrTestError string

func (e ErrTestError) Error() string {
	return string(e)
}

const (
	ErrTest         = ErrTestError("test")
	TestMsg         = "test message"
	FileName        = "error_test.go"
	ExpectedErrText = "InvalidArgument (3): test message"
)

type ErrorTestSuite struct {
	suite.Suite
}

func Test_ErrorTestSuite(t *testing.T) {
	suite.Run(t, new(ErrorTestSuite))
}

func (s *ErrorTestSuite) TestNew() {
	r := s.Require()
	_, _, line := getStackInfo(0)
	err := Wrap(ErrTest, TestMsg, codes.InvalidArgument)
	r.Equal(TestMsg, err.Message)
	r.Equal(codes.InvalidArgument, err.Code)
	r.Equal(FileName, err.File)
	r.Equal("(*ErrorTestSuite).TestNew()", err.Function)
	r.Equal(line+1, err.Line)
}

func (s *ErrorTestSuite) TestNew_genericFunc() {
	r := s.Require()
	line, err := genericFunc("test")
	errgoError := err.(*Error)
	r.Equal(TestMsg, errgoError.Message)
	r.Equal(codes.InvalidArgument, errgoError.Code)
	r.Equal(FileName, errgoError.File)
	r.Equal("genericFunc[...]()", errgoError.Function)
	r.Equal(line, errgoError.Line)
	r.Equal(ExpectedErrText, err.Error())
}

func (s *ErrorTestSuite) TestError() {
	r := s.Require()
	err := Wrap(ErrTest, "test message", codes.InvalidArgument)
	r.Equal(ExpectedErrText, err.Error())
}

func (s *ErrorTestSuite) TestUnwrap() {
	r := s.Require()
	_, _, line := getStackInfo(0)
	err := Wrap(ErrTest, TestMsg, codes.InvalidArgument)
	errgoError := Unwrap(err)
	r.NotNil(errgoError)
	r.Equal(TestMsg, errgoError.Message)
	r.Equal(codes.InvalidArgument, errgoError.Code)
	r.Equal(FileName, errgoError.File)
	r.Equal("(*ErrorTestSuite).TestUnwrap()", errgoError.Function)
	r.Equal(line+1, errgoError.Line)

	errgoError = Unwrap(nil)
	r.Nil(errgoError)

	errgoError = Unwrap(ErrTest)
	r.Nil(errgoError)
}

func (s *ErrorTestSuite) TestUnwrapAll() {
	r := s.Require()
	err := Wrap(ErrTest, TestMsg, codes.InvalidArgument)
	err = Wrap(err, "nested error", codes.Unauthenticated)
	err = Wrap(err, "another nested error", codes.NotFound)

	errors := UnwrapAll(err)
	r.Len(errors, 3)
}

func (s *ErrorTestSuite) TestGetRoot() {
	r := s.Require()
	_, _, line := getStackInfo(0)
	err := Wrap(ErrTest, TestMsg, codes.InvalidArgument)
	err = Wrap(err, "nested error", codes.Unauthenticated)
	err = Wrap(err, "another nested error", codes.NotFound)

	rootErr := GetRoot(err)
	r.Equal(TestMsg, rootErr.Message)
	r.Equal(codes.InvalidArgument, rootErr.Code)
	r.Equal(FileName, rootErr.File)
	r.Equal("(*ErrorTestSuite).TestGetRoot()", rootErr.Function)
	r.Equal(line+1, rootErr.Line)
}

func (s *ErrorTestSuite) TestContainsError() {
	r := s.Require()

	// Test with nil errors
	{
		left := error(nil)
		right := error(nil)

		r.False(ContainsError(left, right))
	}

	// Test with equal errors
	{
		err := errors.New("test error")
		left := err
		right := err

		r.True(ContainsError(left, right))
	}

	// Test with different errors
	{
		left := errors.New("error 1")
		right := errors.New("error 2")

		r.False(ContainsError(left, right))
	}

	// Test with wrapped errors
	{
		inner := ErrTest
		left := Wrap(inner, "outer error", codes.AlreadyExists)
		right := inner

		r.True(ContainsError(left, right))
	}
}

func (s *ErrorTestSuite) TestContainsCode() {
	r := s.Require()
	// Test with nil error
	{
		var err error
		code := codes.NotFound

		r.False(ContainsCode(err, code))
	}

	// Test with matching code
	{
		err := errors.New("test error")
		code := codes.NotFound
		wrapper := Wrap(err, "wrapped error", code)

		r.True(ContainsCode(wrapper, code))
	}

	// Test with non-matching code
	{
		err := errors.New("test error")
		code := codes.Unknown
		wrapper := Wrap(err, "wrapped error", codes.NotFound)

		r.False(ContainsCode(wrapper, code))
	}

	// Test with string
	{
		code := codes.DeadlineExceeded
		err := errors.New(code.String())
		wrapper := Wrap(err, "wrapped error", codes.AlreadyExists)
		str := fmt.Sprint(Wrap(wrapper, "another wrapped error", codes.NotFound).JSON())
		fmt.Println(str)

		r.True(ContainsCode(wrapper, code))
	}

	// Test with unwrapped error
	{
		code := codes.DeadlineExceeded
		err := errors.New(code.String())
		r.True(ContainsCode(err, code))
	}
}

func genericFunc[T any](_ T) (int, error) {
	_, _, line := getStackInfo(0)
	return line + 1, Wrap(ErrTest, TestMsg, codes.InvalidArgument)
}
