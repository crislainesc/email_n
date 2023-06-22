package campaign_test

import (
	"emailn/internal/domain/campaign"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name       = fake.Lorem().Text(10)
	content    = fake.Lorem().Text(20)
	contacts   = []string{fake.Internet().Email()}
	now        = time.Now().Add(-time.Minute)
	created_by = fake.Internet().Email()
	fake       = faker.New()
)

func Test_NewCampaign_CreateCampaign(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(name, content, contacts, created_by)

	assert.NotEqual(campaignEntity, nil)
	assert.Equal(campaignEntity.Name, name)
	assert.Equal(campaignEntity.Content, content)
	assert.Equal(campaignEntity.CreatedBy, created_by)
	assert.Equal(len(campaignEntity.Contacts), len(contacts))
}

func Test_NewCampaign_IDIsNotEmpty(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(name, content, contacts, created_by)

	assert.NotEmpty(campaignEntity.ID)
}

func Test_NewCampaign_StatusIsPending(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(name, content, contacts, created_by)

	assert.Equal(campaignEntity.Status.String(), campaign.Pending.String())
}

func Test_NewCampaign_CreatedOn(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(name, content, contacts, created_by)

	assert.NotNil(campaignEntity.CreatedOn)
	assert.Greater(campaignEntity.CreatedOn, now)
}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, error := campaign.NewCampaign("", content, contacts, created_by)

	assert.Equal("name is required with min 5", error.Error())
}

func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, error := campaign.NewCampaign(fake.Lorem().Text(25), content, contacts, created_by)

	assert.Equal("name is required with max 24", error.Error())
}

func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, error := campaign.NewCampaign(name, "", contacts, created_by)

	assert.Equal("content is required with min 5", error.Error())
}

func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, error := campaign.NewCampaign(name, fake.Lorem().Text(1050), contacts, created_by)

	assert.Equal("content is required with max 1024", error.Error())
}

func Test_NewCampaign_MustValidateContactsMin(t *testing.T) {
	assert := assert.New(t)

	_, error := campaign.NewCampaign(name, content, []string{}, created_by)

	assert.Equal("contacts is required with min 1", error.Error())
}

func Test_NewCampaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, error := campaign.NewCampaign(name, content, []string{"email_invalid"}, created_by)

	assert.Equal("email is invalid", error.Error())
}

func Test_Cancel_ShouldChangeStatusToCanceled(t *testing.T) {
	assert := assert.New(t)
	campaignEntity, _ := campaign.NewCampaign(name, content, contacts, created_by)

	campaignEntity.Cancel()

	assert.Equal(campaignEntity.Status.String(), campaign.Canceled.String())
}

func Test_NewCampaign_MustValidateCreatedBy(t *testing.T) {
	assert := assert.New(t)

	_, error := campaign.NewCampaign(name, content, contacts, "")

	assert.Equal("createdby is invalid", error.Error())
}
