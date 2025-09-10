package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	h "github.com/alfianX/jserver/helper"
	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	status      int
	body        []byte
	wroteHeader bool
	wroteBody   bool
}

func wrapResponseWriter(w gin.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteBody {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if rw.wroteBody {
		return 0, nil
	}
	i, err := rw.ResponseWriter.Write(body)
	if err != nil {
		return 0, err
	}
	rw.body = body
	return i, err
}

func (rw *responseWriter) Body() []byte {
	return rw.body
}

func (s *service) MiddlewareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		defer func() {
			if err := recover(); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}()

		requestBody, err := h.ReadRequestBody(c)
		if err != nil {
			h.Respond(c, err, 0)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		h.RestoreRequestBody(c, requestBody)

		re := regexp.MustCompile(`\r?\n`)
		reqMessage := re.ReplaceAllString(string(requestBody), "")
		reqMessage = strings.ReplaceAll(reqMessage, " ", "")

		logMessage := fmt.Sprintf("path:%s, method: %s,\n requestBody: %v", c.Request.URL.EscapedPath(), c.Request.Method, reqMessage)

		start := time.Now()
		wrapped := wrapResponseWriter(c.Writer)
		c.Writer = wrapped

		c.Next()

		logMessage = fmt.Sprintf("%s,\n respondStatus: %d, respondBody: %s", logMessage, wrapped.Status(), string(wrapped.Body()))
		s.logger.Infof("%s, duration: %v", logMessage, time.Since(start))

	}

}
