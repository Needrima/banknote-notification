package helper

import (
	"github.com/google/uuid"
)

type ErrorBody struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
	Source  string      `json:"source"`
}
type ErrorResponse struct {
	ErrorReference uuid.UUID   `json:"error_reference"`
	TimeStamp      string      `json:"timestamp"`
	Errors         []ErrorBody `json:"errors"`
}
type LogStruct struct {
	TimeStamp       string `json:"@timestamp"`
	Version         string `json:"version"`
	Level           string `json:"level"`
	LevelValue      int    `json:"level_value"`
	StatusCode      string `json:"statusCode"`
	PayLoad         string `json:"pay_load"`
	Message         string `json:"message"`
	LoggerName      string `json:"logger_name"`
	AppName         string `json:"app_name"`
	Path            string `json:"path"`
	Method          string `json:"method"`
	CorrelationId   string `json:"X-Barafiri-Correlation-Id"`
	UserAgent       string `json:"User-Agent"`
	ResponseTime    string `json:"X-Barafiri-Response-Time"`
	ApplicationHost string `json:"X-Barafiri-Application-Host"`
	ForwardedFor    string `json:"X-Forwarded-For"`
}
