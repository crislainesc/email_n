package endpoints_test

import (
	"emailn/internal/endpoints"
	"emailn/internal/internalerrors"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HandlerError_WhenEndpointsReturnsInternalError(t *testing.T) {
	assert := assert.New(t)
	endpoint := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, 0, internalerrors.ErrorInternal
	}
	handlerFunc := endpoints.HandlerError(endpoint)

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	handlerFunc.ServeHTTP(res, req)

	assert.Equal(http.StatusInternalServerError, res.Code)
	assert.Contains(res.Body.String(), internalerrors.ErrorInternal.Error())
}

func Test_HandlerError_WhenEndpointsReturnsDomainError(t *testing.T) {
	assert := assert.New(t)
	endpoint := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, 0, errors.New("domain error")
	}
	handlerFunc := endpoints.HandlerError(endpoint)

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	handlerFunc.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
	assert.Contains(res.Body.String(), "domain error")
}

func Test_HandlerError_WhenEndpointsReturnsObjAndStatus(t *testing.T) {
	assert := assert.New(t)
	type bodyForTest struct {
		Id string
	}
	objExpected := bodyForTest{Id: "2"}
	endpoint := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return objExpected, 201, nil
	}
	handlerFunc := endpoints.HandlerError(endpoint)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	handlerFunc.ServeHTTP(res, req)
	println("res")
	println(res.Code)
	println("res")

	assert.Equal(http.StatusCreated, res.Code)
	objReturned := bodyForTest{}
	json.Unmarshal(res.Body.Bytes(), &objReturned)
	assert.Equal(objExpected, objReturned)
}
