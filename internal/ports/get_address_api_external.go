package ports

import (
	"context"

	"github.com/alef-daniel/challenge-multithreading/internal/application/domain"
)

type GetAddressAPIExternal interface {
	GetAddress(ctx context.Context, cep string) (*domain.Address, error)
}
