package database

import (
	"emailn/internal/domain/campaign"

	"gorm.io/gorm"
)

type CampaignRepository struct {
	Database *gorm.DB
}

func (c *CampaignRepository) Save(campaign *campaign.Campaign) error {
	tx := c.Database.Create(campaign)
	return tx.Error
}

func (c *CampaignRepository) Get() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign

	tx := c.Database.Find(&campaigns)

	return campaigns, tx.Error
}

func (c *CampaignRepository) GetById(id string) (*campaign.Campaign, error) {
	var campaign campaign.Campaign

	tx := c.Database.First(&campaign, id)

	return &campaign, tx.Error
}
