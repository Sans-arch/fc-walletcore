package createtransaction

import (
	"github.com/Sans-arch/fc-walletcore/internal/entity"
	"github.com/Sans-arch/fc-walletcore/internal/gateway"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string
	AccountIDTo string
	Amount float64
}

type CreateTransactionOutputDTO struct {
	ID string
}

type CreateTransactionUsecase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway gateway.AccountGateway
}

func NewTransactionUsecase(transactionGateway gateway.TransactionGateway, accountGateway gateway.AccountGateway) *CreateTransactionUsecase {
	return &CreateTransactionUsecase{
		TransactionGateway: transactionGateway, 
		AccountGateway: accountGateway,
	}
}

func (uc *CreateTransactionUsecase) Execute(input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	accountFrom, err := uc.AccountGateway.Get(input.AccountIDFrom)
	if err != nil {
		return nil, err
	}
	accountTo, err := uc.AccountGateway.Get(input.AccountIDTo)
	if err != nil {
		return nil, err
	}
	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
	if err != nil {
		return nil, err
	}
	err = uc.TransactionGateway.Create(transaction)
	if err != nil {
		return nil, err
	}
	return &CreateTransactionOutputDTO{ID: transaction.ID}, nil
}