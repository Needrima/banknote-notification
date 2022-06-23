package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

var (
	ValidationError = "VALIDATION_ERROR"
	RedisSetupError = "REDIS_SETUP_ERROR"
	NoRecordError   = "NO_RECORD_FOUND_ERROR"
	NoResourceError = "INVALID_RESOURCE_ERROR"
	CreateError     = "CREATE_ERROR"
	UpdateError     = "UPDATE_ERROR"
	//EnvError        = "ENV_ERROR"
	LogError            = "LOG_ERROR"
	MongoDBError        = "MONGO_DB_ERROR"
	ServerError         = "SERVER_ERROR"
	InvalidScheduleDate = "INVALID_SCHEDULE_DATE_ERROR" // You cannot schedule a task in the past. You must provide a future date
)

func (err ErrorResponse) Error() string {
	var errorBody ErrorBody
	return fmt.Sprintf("%v", errorBody)

}
func ErrorArrayToError(errorBody []validator.FieldError) error {
	var errorResponse ErrorResponse
	errorResponse.TimeStamp = time.Now().Format(time.RFC3339)
	errorResponse.ErrorReference = uuid.New()

	for _, value := range errorBody {
		body := ErrorBody{Code: ValidationError, Source: Config.ServiceName, Message: value.Error()}
		errorResponse.Errors = append(errorResponse.Errors, body)
	}
	return errorResponse
}
func ErrorMessage(code string, message string) error {
	var errorResponse ErrorResponse
	errorResponse.TimeStamp = time.Now().Format(time.RFC3339)
	errorResponse.ErrorReference = uuid.New()
	errorResponse.Errors = append(errorResponse.Errors, ErrorBody{Code: code, Source: Config.ServiceName, Message: message})
	return errorResponse
}
