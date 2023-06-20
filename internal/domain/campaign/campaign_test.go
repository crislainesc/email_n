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

	campaign := NewCampaign(name, content, contacts)

	assert.NotEqual(campaign, nil)
	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
}

func TestNewCampaign_IDIsNotEmpty(t *testing.T) {
	assert := assert.New(t)

	campaign := NewCampaign(name, content, contacts)

	assert.NotEmpty(campaign.ID)
}

func TestNewCampaign_CreatedOn(t *testing.T) {
	assert := assert.New(t)

	campaign := NewCampaign(name, content, contacts)

	assert.NotNil(campaign.CreatedOn)
	assert.Greater(campaign.CreatedOn, now)
}
