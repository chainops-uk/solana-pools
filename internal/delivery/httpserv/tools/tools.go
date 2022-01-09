package tools

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type HandlerFunc func(*gin.Context) (interface{}, error)
type WSHandlerFunc func(*gin.Context, []byte) (interface{}, error)

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
					context.Error(err)
				} else {
					context.JSON(t.code, gin.H{
						"error": t.Error(),
					})
				}
			default:
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
				context.Error(err)
			}

			return
		}

		context.JSON(http.StatusOK, result)
	}
}

func WSMust(handlerFunc WSHandlerFunc, TimeSleep time.Duration) gin.HandlerFunc {
	return func(context *gin.Context) {
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		ws, err := upGrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			context.Error(err)
			return
		}
		defer ws.Close()
		ws.SetReadLimit(0)

		rch := make(chan []byte)
		cch := make(chan struct{})
		go func() {

			for {
				_, p, err := ws.ReadMessage()
				if err != nil {
					if _, ok := err.(*websocket.CloseError); ok {
						close(cch)
						break
					}
					context.Error(err)
					break
				}
				rch <- p
			}
			close(rch)
		}()

		for {
			var b []byte
			select {
			case <-cch:
				return
			case message, ok := <-rch:
				if !ok {
					if err := ws.WriteJSON(gin.H{
						"error": "internal server error",
					}); err != nil {
						context.Error(err)
					}
					return
				}
				b = message
			default:
				b = nil
			}

			result, err := handlerFunc(context, b)
			if err != nil {
				switch t := err.(type) {
				case *Status:
					if t.code > 499 {
						if err := ws.WriteJSON(gin.H{
							"error": "internal server error",
						}); err != nil {
							context.Error(err)
						}
						context.Error(err)
					} else {
						if err := ws.WriteJSON(gin.H{
							"error": t.Error(),
						}); err != nil {
							context.Error(err)
						}
					}
				default:
					if err := ws.WriteJSON(gin.H{
						"error": "internal server error",
					}); err != nil {
						context.Error(err)
					}
					context.Error(err)
				}
				return
			}

			if err := ws.WriteJSON(gin.H{
				"data": result,
			}); err != nil {
				context.Error(err)
				return
			}

			time.Sleep(TimeSleep)
		}
	}
}
