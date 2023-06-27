package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	configuration "walls-notification-service/internal/core/helper/configuration-helper"
	errorhelper "walls-notification-service/internal/core/helper/error-helper"
	logger "walls-notification-service/internal/core/helper/log-helper"

	"github.com/gin-gonic/gin"
)

func CreateHeader(writer *logger.BodyLogWriter, key string, value string) interface{} {
	val, ok := writer.Header()[key]
	if !ok {
		writer.Header().Add(key, value)
		return value
	} else {
		return val
	}
}

func TruncateResponse(response string) string {
	count := 0
	truncatedResponse := ""
	for _, char := range response {
		if count >= 200 {
			break
		}
		truncatedResponse += string(char)
		count++
	}
	return truncatedResponse
}

func LogRequest(c *gin.Context) {

	var payload []byte
	var level string = "INFO"

	if c.Request.Body != nil {
		payload, _ = ioutil.ReadAll(c.Request.Body)
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(payload))
	blw := &logger.BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()

	statusCode := c.Writer.Status()
	if statusCode >= 400 {
		level = "ERROR"
	}

	data, err := json.Marshal(&logger.LogStruct{
		Method:          c.Request.Method,
		Level:           level,
		StatusCode:      strconv.Itoa(statusCode),
		Path:            c.Request.URL.String(),
		UserAgent:       c.Request.Header.Get("User-Agent"),
		ForwardedFor:    c.Request.Header.Get("X-Forwarded-For"),
		ResponseTime:    time.Since(time.Now()).String(),
		PayLoad:         string(payload),
		Message:         http.StatusText(statusCode) + " : " + TruncateResponse(blw.Body.String()) + " ... ",
		Version:         "1",
		CorrelationId:   c.Request.Header.Get("X-Correlation-ID"),
		AppName:         configuration.ServiceConfiguration.ServiceName,
		ApplicationHost: c.Request.Host,
		LoggerName:      "",
		TimeStamp:       time.Now().Format(time.RFC3339),
	})
	if err != nil {
		logger.LogEvent("ERROR", errorhelper.ErrorMessage(errorhelper.LogError, err.Error()))
		log.Fatal(err)
	}
	log.Printf("%s\n", data)
}
