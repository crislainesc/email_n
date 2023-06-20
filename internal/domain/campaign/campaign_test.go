package campaign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCampaign(t *testing.T) {
	assert := assert.New(t)
	name := "Campaign Name"
	content := "Campaign Content"
	contacts := []string{"email@example.com", "email@example.com", "email@example.com"}

	campaign := NewCampaign(name, content, contacts)

	assert.NotEqual(campaign, nil)
	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
}
