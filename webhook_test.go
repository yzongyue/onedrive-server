package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, url string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

type WebHookSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *WebHookSuite) SetupTest() {
	suite.router = setupRouter()
}

func (suite *WebHookSuite) TestValidation() {
	w := performRequest(suite.router, "POST", "/webhook?validationToken=123456", nil)
	suite.Equal(200, w.Code)
	suite.Equal("123456", w.Body.String())

}

func TestWebHookSuite(t *testing.T) {
	suite.Run(t, new(WebHookSuite))
}
