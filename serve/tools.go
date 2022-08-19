package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceID 获取 trace id
func TraceID(c *gin.Context) string {
	if rid := c.GetHeader("X-Trace-Id"); rid != "" {
		return rid
	} else if rid := c.GetHeader("X-Request-Id"); rid != "" {
		return rid
	}
	return uuid.New().String()
}

// Success 返回成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
		"traceId": TraceID(c),
	})
	c.Abort()
}

// Error 返回错误
func Error(c *gin.Context, code int, errorCode, errorMessage string) {
	c.JSON(code, gin.H{
		"success":      false,
		"errorCode":    errorCode,
		"errorMessage": errorMessage,
		"traceId":      TraceID(c),
	})
	c.Abort()
}

// Error2 返回错误
func Error2(c *gin.Context, code int, errorCode string, err error) {
	c.JSON(code, gin.H{
		"success":      false,
		"errorCode":    errorCode,
		"errorMessage": err.Error(),
		"traceId":      TraceID(c),
	})
	c.Abort()
}

// ErrorWithData 返回错误
func ErrorWithData(c *gin.Context, data interface{}, code int, errorCode, errorMessage string) {
	c.JSON(code, gin.H{
		"success":      false,
		"data":         data,
		"errorCode":    errorCode,
		"errorMessage": errorMessage,
		"traceId":      TraceID(c),
	})
	c.Abort()
}
