package api

import (
	"net/http"
	"testing"

	"net/http/httptest"
	"github.com/cassavahq/telo/conf"
	"github.com/stretchr/testify/assert"
)

func TestCanSayHello(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := httptest.NewRequest("GET", "/", nil)
	rsp := httptest.NewRecorder()

	a.handler.NewContext(req, rsp)
	assert.Equal(t, http.StatusOK, rsp.Code)
}

func TestNoAuthProvided(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := httptest.NewRequest("GET", "/private", nil)
	rsp := httptest.NewRecorder()
 a.handler.NewContext(req, rsp)


	assert.Equal(t, http.StatusBadRequest, rsp.Code)
}

func TestBadAuthProvided(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := httptest.NewRequest("GET", "/private", nil)
	req.Header.Add("Authorization", "Bearer nonsense")
	rsp := httptest.NewRecorder()
	a.handler.NewContext(req, rsp)


	assert.Equal(t, http.StatusUnauthorized, rsp.Code)
}

func TestAuthOk(t *testing.T) {
	config := &conf.Config{
		JWTSecret: "secret",
	}
	a := NewAPI(config)

	req := httptest.NewRequest("GET", "/private", nil)
	req.Header.Add("Authorization", "Bearer nonsense")
	rsp := httptest.NewRecorder()
	a.handler.NewContext(req, rsp)

	assert.Equal(t, http.StatusUnauthorized, rsp.Code)
}