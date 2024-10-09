package gateway

import "github.com/Sans-arch/fc-walletcore/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
