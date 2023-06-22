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

func (r *repositoryMock) Get() ([]Campaign, error) {
	args := r.Called()

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]Campaign), nil
}

func (r *repositoryMock) GetById(id string) (*Campaign, error) {
	args := r.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Campaign), nil
}

func (r *repositoryMock) Update(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *repositoryMock) Delete(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

var (
	newCampaign = contract.NewCampaignInput{
		Name:    "TestCreateCampaign",
		Content: "test content",
		Emails:  []string{"email@example.com"},
	}
	service           = ServiceImp{}
	errSomethingWrong = errors.New("something went wrong")
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("Save", mock.Anything).Return(nil)
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

func Test_GetById_ReturnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.MatchedBy(func(id string) bool {
		return id == campaign.ID
	})).Return(campaign, nil)
	service.Repository = repositoryMock

	campaignReturned, _ := service.GetById(campaign.ID)

	assert.Equal(campaign.ID, campaignReturned.ID)
	assert.Equal(campaign.Name, campaignReturned.Name)
	assert.Equal(campaign.Content, campaignReturned.Content)
	assert.Equal(campaign.Status.String(), campaignReturned.Status)
}

func Test_GetById_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(nil, errSomethingWrong)
	service.Repository = repositoryMock

	_, err := service.GetById(campaign.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Cancel_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(nil, errSomethingWrong)
	repositoryMock.On("Update", mock.Anything).Return(errSomethingWrong)
	service.Repository = repositoryMock

	err := service.Cancel(campaign.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Cancel_ReturnErrorIfStatusIsNotPending(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	errorExpected := errors.New("campaign is not pending")
	repositoryMock.On("GetById", mock.Anything).Return(&Campaign{}, nil)
	repositoryMock.On("Update", mock.Anything).Return(errorExpected)
	service.Repository = repositoryMock

	err := service.Cancel(campaign.ID)

	assert.Equal(err.Error(), errorExpected.Error())
}

func Test_Cancel_ReturnErrorWhenUpdateFail(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaign, nil)
	repositoryMock.On("Update", mock.Anything).Return(errSomethingWrong)
	service.Repository = repositoryMock

	err := service.Cancel(campaign.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Cancel_ShouldReturnNilIfSuccess(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaign, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	err := service.Cancel(campaign.ID)

	assert.Nil(err)
}

func Test_Delete_ReturnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(nil, errSomethingWrong)
	repositoryMock.On("Update", mock.Anything).Return(errSomethingWrong)
	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Delete_ReturnErrorIfStatusIsNotPending(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	errorExpected := errors.New("campaign is not pending")
	repositoryMock.On("GetById", mock.Anything).Return(&Campaign{}, nil)
	repositoryMock.On("Update", mock.Anything).Return(errorExpected)
	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Equal(err.Error(), errorExpected.Error())
}

func Test_Delete_ReturnErrorWhenUpdateFail(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaign, nil)
	repositoryMock.On("Update", mock.Anything).Return(errSomethingWrong)
	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Equal(err.Error(), internalerrors.ErrorInternal.Error())
}

func Test_Delete_ShouldReturnNilIfSuccess(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	repositoryMock := new(repositoryMock)
	repositoryMock.On("GetById", mock.Anything).Return(campaign, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)
	service.Repository = repositoryMock

	err := service.Delete(campaign.ID)

	assert.Nil(err)
}
