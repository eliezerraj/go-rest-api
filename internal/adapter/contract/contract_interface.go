package contract

import(
	"context"
	"github.com/go-rest-api/internal/model"
)

type BalanceServiceAdapterPort interface {
	ListBalance() ([]model.Balance ,error)
	ListBalanceById(pk string, sk string) ([]model.Balance ,error)
	GetBalance(account string) (model.Balance ,error)
	AddBalance(balance model.Balance) (model.Balance, error)
}

type BalanceRepositoryAdapterPort interface {
	ListBalance(ctx context.Context) ([]model.Balance ,error)
	ListBalanceById(ctx context.Context, pk string, sk string) ([]model.Balance ,error)
	GetBalance(ctx context.Context, account string) (model.Balance ,error)
	AddBalance(ctx context.Context, balance model.Balance) (model.Balance, error)
}

type MetricsServiceAdapterPort interface {
	Health() (bool)
	StressCPU(count int) string
}
