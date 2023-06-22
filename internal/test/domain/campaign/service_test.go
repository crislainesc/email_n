package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	"emailn/internal/internalerrors"
	"emailn/internal/test/internalmock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	newCampaign = contract.NewCampaignInput{
		Name:    "TestCreateCampaign",
		Content: "test content",
		Emails:  []string{"email@example.com"},
	}
	service           = campaign.ServiceImp{}
	errSomethingWrong = errors.New("something went wrong")
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.False(errors.Is(internalerrors.ErrorInternal, err))
}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)

	_, err := service.Create(contract.NewCampaignInput{})

	assert.NotNil(err)
	assert.Equal("name is required with min 5", err.Error())
}

func Test_Create_CreateCampaign(t *testing.T) {
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name || campaign.Content != newCampaign.Content || len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)
	service.Repository = repositoryMock
	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositoryCreate(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to create on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrorInternal, err))
}

func Test_GetById_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.MatchedBy(func(id string) bool {
		return id == campaignEntity.ID
	})).Return(campaignEntity, nil)
	service.Repository = repositoryMock

	campaignReturned, _ := service.GetById(campaignEntity.ID)

	assert.Equal(campaignEntity.ID, campaignReturned.ID)
	assert.Equal(campaignEntity.Name, campaignReturned.Name)
	assert.Equal(campaignEntity.Content, campaignReturned.Content)
	assert.Equal(campaignEntity.Status.String(), campaignReturned.Status)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(nil, errSomethingWrong)
	service.Repository = repositoryMock

	_, err := service.GetById(campaignEntity.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Delete_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(nil, errSomethingWrong)
	repositoryMock.On("Delete", mock.Anything).Return(errSomethingWrong)
	service.Repository = repositoryMock

	err := service.Delete(campaignEntity.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Delete_ReturnErrorIfStatusIsNotPending(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	errorExpected := errors.New("campaign is not pending")
	repositoryMock.On("GetById", mock.Anything).Return(&campaign.Campaign{}, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errorExpected)
	service.Repository = repositoryMock

	err := service.Delete(campaignEntity.ID)

	assert.Equal(err.Error(), errorExpected.Error())
}

func Test_Delete_ReturnErrorWhenUpdateFail(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaignEntity, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errSomethingWrong)
	service.Repository = repositoryMock

	err := service.Delete(campaignEntity.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Delete_ShouldReturnNilIfSuccess(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaignEntity, nil)
	repositoryMock.On("Delete", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	err := service.Delete(campaignEntity.ID)

	assert.Nil(err)
}
