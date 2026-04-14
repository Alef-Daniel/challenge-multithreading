package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alef-daniel/challenge-multithreading/internal/application/domain"
	"github.com/alef-daniel/challenge-multithreading/pkg"
)

var (
	ErrCepIsEmpty      = errors.New("cep is empty")
	ErrCepNotFound     = errors.New("cep not found")
	ErrViaCepAPI       = errors.New("ViaCEP api error")
	ErrInvalidTypeData = errors.New("invalid type data")
)

type GetAddressViaCepUseCase struct {
	client pkg.Client
}

func (g *GetAddressViaCepUseCase) GetAddress(ctx context.Context, cep string) (*domain.Address, error) {
	if cep == "" {
		return nil, ErrCepIsEmpty
	}

	url := g.BuildURL(ctx, cep)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error make request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 0 && resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, ErrCepNotFound
		default:
			return nil, ErrViaCepAPI
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error read response body: %w", err)
	}

	address, err := g.BuildResponse(ctx, body)
	if err != nil {
		return nil, err
	}

	return address, nil

}

func (g *GetAddressViaCepUseCase) BuildURL(ctx context.Context, cep string) string {
	return fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
}

func (g *GetAddressViaCepUseCase) BuildResponse(ctx context.Context, response []byte) (*domain.Address, error) {
	responseMap := make(map[string]interface{})
	if len(response) == 0 {
		return nil, ErrViaCepAPI
	}

	address := &domain.Address{}
	err := json.Unmarshal(response, &responseMap)
	if err != nil {
		return nil, err
	}

	if responseMap["cep"] != nil {
		cep, ok := responseMap["cep"].(string)
		if !ok {
			return nil, ErrInvalidTypeData
		}
		address.Cep = cep
	}

	if responseMap["logradouro"] != nil {
		logradouro, ok := responseMap["logradouro"].(string)
		if !ok {
			return nil, ErrInvalidTypeData
		}
		address.Logradouro = logradouro
	}

	if responseMap["uf"] != nil {
		uf, ok := responseMap["uf"].(string)
		if !ok {
			return nil, ErrInvalidTypeData
		}

		address.UF = uf
	}

	if responseMap["Bairro"] != nil {
		bairro, ok := responseMap["Bairro"].(string)
		if !ok {
			return nil, ErrInvalidTypeData
		}
		address.Bairro = bairro
	}

	address.Provider = domain.ProviderViaCEP

	return address, nil

}

func NewGetAddressViaCepUseCase(client pkg.Client) *GetAddressViaCepUseCase {
	return &GetAddressViaCepUseCase{client: client}
}
