package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"notification-service/internal/core/helper"
	"strconv"
	"time"
)

func CreateHeader(writer *helper.BodyLogWriter, key string, value string) interface{} {
	val, ok := writer.Header()[key]
	if !ok {
		writer.Header().Add(key, value)
		return value
	} else {
		return val
	}
}

func LogRequest(c *gin.Context) {
	blw := &helper.BodyLogWriter{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	correlationId := uuid.New().String()
	CreateHeader(blw, "x-correlation-id", correlationId)
	c.Writer = blw
	c.Next()
	statusCode := c.Writer.Status()
	level := "INFO"
	if statusCode >= 400 {
		level = "ERROR"
	}

	count := 0
	truncatedResponse := ""
	for _, char := range blw.Body.String() {
		if count >= 200 {
			break
		}
		truncatedResponse += string(char)
		count++
	}
	payload, _ := c.GetRawData()
	data, err := json.Marshal(&helper.LogStruct{
		Method:          c.Request.Method,
		Level:           level,
		StatusCode:      strconv.Itoa(statusCode),
		Path:            c.Request.URL.String(),
		UserAgent:       c.Request.Header.Get("User-Agent"),
		ForwardedFor:    c.Request.Header.Get("X-Forwarded-For"),
		ResponseTime:    time.Since(time.Now()).String(),
		PayLoad:         string(payload),
		Message:         http.StatusText(statusCode) + " : " + truncatedResponse + " ... ",
		Version:         "1",
		CorrelationId:   correlationId,
		AppName:         helper.Config.ServiceName,
		ApplicationHost: c.Request.Host,
		LoggerName:      "",
		TimeStamp:       time.Now().Format(time.RFC3339),
	})
	if err != nil {
		helper.LogEvent("ERROR", helper.ErrorMessage(helper.LogError, err.Error()))
		log.Fatal(err)
	}
	log.Printf("%s\n", data)
}
