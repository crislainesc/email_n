package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/errors"
	"errors"
)

type Service interface {
	Create(newCampaign contract.NewCampaignInput) (string, error)
	GetById(id string) (*contract.GetCampaignByIdOutput, error)
	Cancel(id string) error
	Delete(id string) error
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
		return nil, internalerrors.ProcessErrorToReturn(err)
	}

	return &contract.GetCampaignByIdOutput{
		ID:                   campaign.ID,
		Name:                 campaign.Name,
		Content:              campaign.Content,
		Status:               campaign.Status.String(),
		AmountOfEmailsToSend: len(campaign.Contacts),
	}, nil
}

func (s *ServiceImp) Cancel(id string) error {

	campaign, err := s.Repository.GetById(id)

	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign.Status != Pending {
		return errors.New("campaign is not pending")
	}

	campaign.Cancel()
	err = s.Repository.Update(campaign)

	if err != nil {
		return internalerrors.ErrorInternal
	}

	return nil
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
