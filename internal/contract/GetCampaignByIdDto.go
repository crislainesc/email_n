package contract

type GetCampaignByIdOutput struct {
	ID                   string
	Name                 string
	Content              string
	Status               string
	AmountOfEmailsToSend int
	CreatedBy            string
}
