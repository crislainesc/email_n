package endpoints

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalmock "emailn/internal/test/mock"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignsGetById_ShouldReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaignResponse := contract.GetCampaignByIdOutput{
		ID:      "1",
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(20),
		Status:  campaign.Pending.String(),
	}
	service := new(internalmock.CampaignServiceMock)
	service.On("GetById", mock.Anything).Return(&campaignResponse, nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	response, status, _ := handler.CampaignGetById(res, req)

	assert.Equal(http.StatusOK, status)
	assert.Equal(campaignResponse.ID, response.(*contract.GetCampaignByIdOutput).ID)
	assert.Equal(campaignResponse.Name, response.(*contract.GetCampaignByIdOutput).Name)
}

func Test_CampaignsGetById_WhenSomethingWrong(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	errorExpected := errors.New("something wrong")
	service.On("GetById", mock.Anything).Return(nil, errorExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	_, status, err := handler.CampaignGetById(res, req)

	assert.Equal(http.StatusNotFound, status)
	assert.Equal(err.Error(), errorExpected.Error())
}
