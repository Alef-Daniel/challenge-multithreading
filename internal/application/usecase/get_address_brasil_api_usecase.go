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
	ErrBrasilAPI = errors.New("Brasil api error")
)

type GetAddressBrasilAPIUseCase struct {
	client pkg.Client
}

func (g *GetAddressBrasilAPIUseCase) GetAddress(ctx context.Context, cep string) (*domain.Address, error) {
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
			return nil, ErrBrasilAPI
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

func (g *GetAddressBrasilAPIUseCase) BuildURL(ctx context.Context, cep string) string {
	return fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
}

func (g *GetAddressBrasilAPIUseCase) BuildResponse(ctx context.Context, response []byte) (*domain.Address, error) {
	responseMap := make(map[string]interface{})
	if len(response) == 0 {
		return nil, ErrBrasilAPI
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

	if responseMap["street"] != nil {
		logradouro, ok := responseMap["street"].(string)
		if !ok {
			return nil, ErrInvalidTypeData
		}
		address.Logradouro = logradouro
	}

	if responseMap["state"] != nil {
		uf, ok := responseMap["state"].(string)
		if !ok {
			return nil, ErrInvalidTypeData
		}

		address.UF = uf
	}

	if responseMap["neighborhood"] != nil {
		bairro, ok := responseMap["neighborhood"].(string)
		if !ok {
			return nil, ErrInvalidTypeData
		}
		address.Bairro = bairro
	}

	address.Provider = domain.ProviderBrasilAPI

	return address, nil

}

func NewGetAddressBrasilAPIUseCase(client pkg.Client) *GetAddressBrasilAPIUseCase {
	return &GetAddressBrasilAPIUseCase{client: client}
}
