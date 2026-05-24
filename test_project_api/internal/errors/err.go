package errors

type AppError struct {
    Code    int
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

var (
    ErrDepartmentNotFound = &AppError{Code: 404, Message: "department not found"}
    ErrEmployeeNotFound   = &AppError{Code: 404, Message: "employee not found"}
    ErrCycleDetected      = &AppError{Code: 409, Message: "cycle detected in department tree"}
    ErrInvalidParent      = &AppError{Code: 400, Message: "invalid parent_id"}
    ErrDuplicateName      = &AppError{Code: 409, Message: "department name must be unique within same parent"}
    ErrInvalidInput       = &AppError{Code: 400, Message: "invalid input"}
)