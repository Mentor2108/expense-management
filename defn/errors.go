package defn

import "errors"

const (
	ErrCodeFailedToParseRequestBody      = "failed-to-parse-request-body"
	ErrCodeMissingRequiredField          = "missing-required-field"
	ErrCodeInputInvalidFormat            = "invalid-input-field-format"
	ErrCodeLoginFailed                   = "login-failed"
	ErrCodeNoDataFound                   = "no-data-found"
	ErrCodeInvalidToken                  = "invalid-token"
	ErrCodeDatabaseCreateOperationFailed = "database-create-operation-failed"
	ErrCodeDatabaseUpdateOperationFailed = "database-update-operation-failed"
	ErrCodeDatabaseGetOperationFailed    = "database-get-operation-failed"
	ErrCodeDatabaseDeleteOperationFailed = "database-delete-operation-failed"
	ErrCodeUnexpectedError               = "unexpected-error"
)

var (
	ErrFailedToParseRequestBody      = errors.New("failed to parse request body: {error}")
	ErrMissingRequiredField          = errors.New("missing required field: {field}")
	ErrInputInvalidFormat            = errors.New("invalid format for field: {message}")
	ErrLoginFailed                   = errors.New("invalid email or password")
	ErrNoDataFound                   = errors.New("no data found in database")
	ErrInvalidToken                  = errors.New("invalid token provided")
	ErrDatabaseCreateOperationFailed = errors.New("database create action failed: {error}")
	ErrDatabaseUpdateOperationFailed = errors.New("database update action failed: {error}")
	ErrDatabaseGetOperationFailed    = errors.New("database get action failed: {error}")
	ErrDatabaseDeleteOperationFailed = errors.New("database delete action failed: {error}")
	ErrUnexpectedError               = errors.New("unexpected error occurred: {error}")
)
