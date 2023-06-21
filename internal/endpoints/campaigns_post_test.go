package endpoints

import (
	"bytes"
	"emailn/internal/contract"
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
	body = contract.NewCampaign{
		Name:    fake.Lorem().Text(10),
		Content: fake.Lorem().Text(20),
		Emails:  []string{fake.Internet().Email(), fake.Internet().Email()},
	}
)

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) Create(newCampaign contract.NewCampaign) (string, error) {
	args := s.Called(newCampaign)
	return args.String(0), args.Error(1)
}

func Test_CampaignsPost_ShouldSaveNewCampaign(t *testing.T) {
	assert := assert.New(t)
	service := &serviceMock{}
	service.On("Create", mock.MatchedBy(func(request contract.NewCampaign) bool {
		if request.Name == body.Name && request.Content == body.Content && len(request.Emails) == len(body.Emails) {
			return true
		} else {
			return false
		}
	})).Return("1", nil)
	handler := Handler{CampaignService: service}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("POST", "/", &buf)
	res := httptest.NewRecorder()

	_, status, err := handler.CampaignPost(res, req)

	assert.Equal(201, status)
	assert.Nil(err)
}

func Test_CampaignsPost_ShouldInformError(t *testing.T) {
	assert := assert.New(t)

	service := &serviceMock{}
	service.On("Create", mock.Anything).Return("", fmt.Errorf("error"))
	handler := Handler{CampaignService: service}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("POST", "/", &buf)
	res := httptest.NewRecorder()

	_, status, err := handler.CampaignPost(res, req)

	assert.Equal(400, status)
	assert.NotNil(err)
}
