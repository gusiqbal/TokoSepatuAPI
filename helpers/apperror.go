package helpers

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(code int, messaage string) *AppError {
	return &AppError{
		Code:    code,
		Message: messaage,
	}
}
