package create_account

import (
	"github.com/Sans-arch/fc-walletcore/internal/entity"
	"github.com/Sans-arch/fc-walletcore/internal/gateway"
)

type CreateAccountInputDTO struct {
	ClientID string `json:"client_id"`
}

type CreateAccountOutputDTO struct {
	ID string
}

type CreateAccountUsecase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUsecase(a gateway.AccountGateway, c gateway.ClientGateway) *CreateAccountUsecase {
	return &CreateAccountUsecase{AccountGateway: a, ClientGateway: c}
}

func (uc *CreateAccountUsecase) Execute(input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	client, err := uc.ClientGateway.Get(input.ClientID)
	if err != nil {
		return nil, err
	}
	account := entity.NewAccount(client)
	err = uc.AccountGateway.Save(account)
	if err != nil {
		return nil, err
	}
	return &CreateAccountOutputDTO{ID: account.ID}, nil
}
