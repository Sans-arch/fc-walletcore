package gateway

import "github.com/Sans-arch/fc-walletcore/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	Get(id string) (*entity.Account, error)
}