package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/errors"
)

type Service interface {
	Create(newCampaign contract.NewCampaignInput) (string, error)
	GetById(id string) (*contract.GetCampaignByIdOutput, error)
}

type ServiceImp struct {
	Repository Repository
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaignInput) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return "", err
	}

	err = s.Repository.Save(campaign)

	if err != nil {
		return "", internalerrors.ErrorInternal
	}

	return campaign.ID, nil
}

func (s *ServiceImp) GetById(id string) (*contract.GetCampaignByIdOutput, error) {

	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return nil, internalerrors.ErrorInternal
	}

	return &contract.GetCampaignByIdOutput{
		ID:      campaign.ID,
		Name:    campaign.Name,
		Content: campaign.Content,
		Status:  campaign.Status.String(),
	}, nil
}
