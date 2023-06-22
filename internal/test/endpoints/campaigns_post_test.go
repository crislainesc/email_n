package endpoints_test

import (
	"bytes"
	"context"
	"emailn/internal/contract"
	"emailn/internal/endpoints"
	"emailn/internal/test/internalmock"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	fake = faker.New()
	body = contract.NewCampaignInput{
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(20),
		Emails:  []string{fake.Internet().Email(), fake.Internet().Email()},
	}
	createdByExpected = fake.Internet().Email()
)

func setup(body contract.NewCampaignInput, createdByExpected string) (*http.Request, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("POST", "/", &buf)

	ctx := context.WithValue(req.Context(), endpoints.EmailKey, createdByExpected)
	req = req.WithContext(ctx)
	res := httptest.NewRecorder()

	return req, res
}

func Test_CampaignsPost_ShouldCreateNewCampaign(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.MatchedBy(func(request contract.NewCampaignInput) bool {
		if request.Name == body.Name &&
			request.Content == body.Content &&
			len(request.Emails) == len(body.Emails) &&
			request.CreatedBy == createdByExpected {
			return true
		} else {
			return false
		}
	})).Return("1", nil)
	handler := endpoints.Handler{CampaignService: service}
	req, res := setup(body, createdByExpected)

	_, status, err := handler.CampaignPost(res, req)

	assert.Equal(http.StatusCreated, status)
	assert.Nil(err)
}

func Test_CampaignsPost_ShouldInformError(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	handler := endpoints.Handler{CampaignService: service}
	req, res := setup(body, createdByExpected)

	_, _, err := handler.CampaignPost(res, req)

	assert.NotNil(err)
}
