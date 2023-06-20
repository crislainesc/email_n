package campaign

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	name     = "Campaign Name"
	content  = "Campaign Content"
	contacts = []string{"email@example.com", "email@example.com", "email@example.com"}
	now      = time.Now().Add(-time.Minute)
)

func TestNewCampaign_CreateCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.NotEqual(campaign, nil)
	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
}

func TestNewCampaign_IDIsNotEmpty(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.NotEmpty(campaign.ID)
}

func TestNewCampaign_CreatedOn(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.NotNil(campaign.CreatedOn)
	assert.Greater(campaign.CreatedOn, now)
}

func TestNewCampaign_MustValidateName(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign("", content, contacts)

	assert.Equal("name is required", error.Error())
}

func TestNewCampaign_MustValidateContent(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, "", contacts)

	assert.Equal("content is required", error.Error())
}

func TestNewCampaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, error := NewCampaign(name, content, []string{})

	assert.Equal("contacts is required", error.Error())
}
