package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/errors"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Save(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *repositoryMock) Create(newCampaign contract.NewCampaign) (string, error) {
	args := r.Called(newCampaign)
	return args.String(0), args.Error(0)
}

var (
	newCampaign = contract.NewCampaign{
		Name:    "TestCreateCampaign",
		Content: "test content",
		Emails:  []string{"email@example.com"},
	}
	service = Service{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("Save", mock.Anything).Return(nil)
	repositoryMock.On("Create", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	id, err := service.Create(newCampaign)

	assert.NotEmpty(id)
	assert.Nil(err)
}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)
	newCampaign.Name = ""

	_, err := service.Create(newCampaign)

	assert.NotNil(err)
	assert.Equal("name is required", err.Error())
}

func Test_Create_SaveCampaign(t *testing.T) {
	repositoryMock := new(repositoryMock)
	repositoryMock.On("Save", mock.MatchedBy(func(campaign *Campaign) bool {
		if campaign.Name != newCampaign.Name || campaign.Content != newCampaign.Content || len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)
	service.Repository = repositoryMock
	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("Save", mock.Anything).Return(errors.New("error to save on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrorInternal, err))
}
