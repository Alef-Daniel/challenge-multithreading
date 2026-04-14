package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/alef-daniel/challenge-multithreading/internal/application/domain"
	"github.com/alef-daniel/challenge-multithreading/internal/ports"
)

type ProcessAddressUseCase struct {
	viaCEP    ports.GetAddressAPIExternal
	BrasilAPI ports.GetAddressAPIExternal
}

func NewProcessAddressUseCase(viaCEP, brasilAPI ports.GetAddressAPIExternal) *ProcessAddressUseCase {
	return &ProcessAddressUseCase{
		viaCEP:    viaCEP,
		BrasilAPI: brasilAPI,
	}
}

func (p *ProcessAddressUseCase) Execute(ctx context.Context, cep string) (*domain.Address, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	type AddressResult struct {
		Address *domain.Address
		Err     error
	}
	ch := make(chan *AddressResult, 2)

	if cep == "" {
		return nil, ErrCepNotFound
	}

	go func() {
		address, err := p.viaCEP.GetAddress(ctx, cep)
		select {
		case ch <- &AddressResult{
			Address: address,
			Err:     err,
		}:
		case <-ctx.Done():
		}

	}()

	go func() {
		address, err := p.BrasilAPI.GetAddress(ctx, cep)
		select {
		case ch <- &AddressResult{
			Address: address,
			Err:     err,
		}:
		case <-ctx.Done():
		}
	}()

	select {
	case result := <-ch:
		return result.Address, result.Err

	case <-ctx.Done():
		return nil, fmt.Errorf("Timeout ")
	}

}
