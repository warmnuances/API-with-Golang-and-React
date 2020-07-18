package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/orders/all", nil)
	if err!= nil {
		fmt.Println(err)
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

