package mock

import (
	"emailn/internal/contract"

	"github.com/stretchr/testify/mock"
)

type CampaignServiceMock struct {
	mock.Mock
}

func (s *CampaignServiceMock) Create(newCampaign contract.NewCampaignInput) (string, error) {
	args := s.Called(newCampaign)
	return args.String(0), args.Error(1)
}

func (s *CampaignServiceMock) GetById(id string) (*contract.GetCampaignByIdOutput, error) {
	args := s.Called(id)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*contract.GetCampaignByIdOutput), nil
}

func (s *CampaignServiceMock) Cancel(id string) error {
	args := s.Called(id)

	println(args.Error(0))

	return args.Error(0)
}
