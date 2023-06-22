package endpoints

import (
	internalmock "emailn/internal/test/mock"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignDelete_ShouldReturnTrueIfSuccess(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	service.On("Delete", mock.Anything).Return(nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("PATCH", "/campaigns/1", nil)
	res := httptest.NewRecorder()

	response, status, _ := handler.CampaignDelete(res, req)

	assert.Equal(200, status)
	assert.Equal(response.(bool), true)
}

func Test_CampaignDelete_ShouldReturnErrorIfSomethingWrong(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	errorExpected := errors.New("something wrong")
	service.On("Delete", mock.Anything).Return(errorExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("PATCH", "/campaigns/1", nil)
	res := httptest.NewRecorder()

	_, status, err := handler.CampaignDelete(res, req)

	assert.Equal(400, status)
	assert.Equal(err.Error(), errorExpected.Error())
}
