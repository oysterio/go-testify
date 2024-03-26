package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Request check
	newRequire := require.New(t)
	newRequire.Equal(http.MethodGet, req.Method, "Unexpected request method")
	newRequire.Equal("/cafe", req.URL.Path, "Unexpected request url")

	query := req.URL.Query()
	reqCount, err := strconv.Atoi(query.Get("count"))
	newAssert := assert.New(t)
	// String to Int conversition check
	newAssert.NoError(err, "Error when reading response")
	reqCity := query.Get("city")
	totalCount := 4
	city := "moscow"

	newRequire.NotNil(reqCount)
	newRequire.NotNil(reqCity)
	newRequire.Equal(totalCount, reqCount, "Unexpected cafe count at the city")
	newRequire.Equal(city, reqCity, "Unexpected city name")

	// Response check
	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status code is 200")
	result := responseRecorder.Result()
	require.NotEmpty(t, result.Body, "Result is empty")
}

func TestMainHandlerWhereIsTheWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=City", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)
	// Status code check
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status code is 400")
	// Response reader error check
	result := responseRecorder.Result()
	body, err := io.ReadAll(result.Body)
	newAssert := assert.New(t)
	newAssert.NoError(err, "Error when reading response")
	// Error message check
	responseBody := "wrong city value"
	assert.Equal(t, string(body), responseBody, "Unexpected response body")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=1000&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)
	handler.ServeHTTP(responseRecorder, req)

	result := responseRecorder.Result()
	body, err := io.ReadAll(result.Body)
	newAssert := assert.New(t)
	// Response reader error check
	newAssert.NoError(err, "Error when reading response")
	list := strings.Split(string(body), ",")
	// Count more than total check
	assert.Len(t, list, totalCount, "Unexpected cafe list size")
}
