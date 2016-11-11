package httputil

import (
	"encoding/json"

	"github.com/MISingularity/deepdash/pkg/log"

	"github.com/gin-gonic/gin"
)

type HTTPInterfaces interface {
	WriteHTTPError(err HTTPError, errMsg string)
	WriteData(statusCode int, data interface{})
	Redirect(statusCode int, redirectURL string)
}

type HTTPError struct {
	// HTTP status code to write into HTTP header.
	// This field should not be marshaled into response
	// JSON body.
	StatusCode int `json:"-"`
	// Deepdash specific error code.
	// Error code can be pragmatically consumed.
	Code int `json:"code"`
	// The error message of the error code.
	// Error message can be printed out and consumed by human.
	Message string `json:"message"`
	Fatal   bool   `json:"-"`
}

type Ginframework struct {
	ginContext *gin.Context
}

func NewGinframework(c *gin.Context) *Ginframework {
	return &Ginframework{c}
}

func (he HTTPError) Error() string {
	b, err := json.Marshal(&he)
	if err != nil {
		panic("unexpected json marshal error")
	}
	return string(b)
}

func (c *Ginframework) WriteHTTPError(err HTTPError, errMsg string) {
	if errMsg != "" {
		log.Errorf("Handler failed! Error Msg = %s", errMsg)
	}
	if err.Fatal {
		log.Fatalf("Fatal Error! Err Msg=%v", err.Message)
	}
	c.ginContext.JSON(err.StatusCode, err)
}

func (c *Ginframework) Redirect(statusCode int, redirectURL string) {
	c.ginContext.Redirect(statusCode, redirectURL)
}

func (c *Ginframework) WriteData(statusCode int, data interface{}) {
	c.ginContext.JSON(statusCode, data)
}
