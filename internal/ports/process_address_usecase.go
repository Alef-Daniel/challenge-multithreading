package ports

import (
	"context"

	"github.com/alef-daniel/challenge-multithreading/internal/application/domain"
)

type ProcessAddressUseCase interface {
	Execute(ctx context.Context, cep string) (*domain.Address, error)
}
