package util

type CodeError struct {
	Code int
	Msg  string
}

func NewCodeError(code int, msg string) error {
	return &CodeError{
		Code: code,
		Msg:  msg,
	}
}

func (codeError *CodeError) Error() string {
	return codeError.Msg
}
