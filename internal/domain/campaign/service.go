package campaign

import (
	"emailn/internal/contract"
	"emailn/internal/internalerrors"
	"errors"
)

type Service interface {
	Create(newCampaign contract.NewCampaignInput) (string, error)
	GetById(id string) (*contract.GetCampaignByIdOutput, error)
	Delete(id string) error
}

type ServiceImp struct {
	Repository Repository
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaignInput) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)

	if err != nil {
		return "", err
	}

	err = s.Repository.Create(campaign)

	if err != nil {
		return "", internalerrors.ErrorInternal
	}

	return campaign.ID, nil
}

func (s *ServiceImp) GetById(id string) (*contract.GetCampaignByIdOutput, error) {

	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return nil, internalerrors.ProcessErrorToReturn(err)
	}

	return &contract.GetCampaignByIdOutput{
		ID:                   campaign.ID,
		Name:                 campaign.Name,
		Content:              campaign.Content,
		Status:               campaign.Status.String(),
		AmountOfEmailsToSend: len(campaign.Contacts),
		CreatedBy:            campaign.CreatedBy,
	}, nil
}

func (s *ServiceImp) Delete(id string) error {
	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign.Status != Pending {
		return errors.New("campaign is not pending")
	}

	err = s.Repository.Delete(campaign)

	if err != nil {
		return internalerrors.ErrorInternal
	}

	return nil
}
