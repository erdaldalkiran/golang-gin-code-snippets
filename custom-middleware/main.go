package main

import (
	"bytes"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: new(bytes.Buffer), ResponseWriter: c.Writer}
		c.Writer = blw
		t := time.Now()
		c.Next()

		// access the status we are sending
		status := c.Writer.Status()
		latency := time.Since(t)
		log.Println("|", status, "|", latency, "|", c.Request.RemoteAddr, "|", c.Request.UserAgent(), "|", c.Request.Method, "|", c.Request.URL, "|", blw.body.String())
	}
}

func main() {
	router := gin.New()
	router.Use(logger())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}
