package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc func(*gin.Context) (interface{}, error)
type Status struct {
	code  int
	error error
}

func NewStatus(code int, error error) *Status {
	return &Status{code: code, error: error}
}

func (s *Status) Error() string {
	return s.error.Error()
}

func Must(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		result, err := handlerFunc(context)
		if err != nil {
			switch t := err.(type) {
			case *Status:
				if t.code > 499 {
					context.JSON(http.StatusInternalServerError, gin.H{
						"error": "internal server error",
					})
				} else {
					context.JSON(t.code, gin.H{
						"error": t.Error(),
					})
				}
			default:
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	}
}
