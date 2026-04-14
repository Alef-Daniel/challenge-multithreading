package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alef-daniel/challenge-multithreading/internal/application/usecase"
	"github.com/alef-daniel/challenge-multithreading/pkg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	httpClient := pkg.NewClient(time.Minute)
	ViaCep := usecase.NewGetAddressViaCepUseCase(*httpClient)
	BrasilAPI := usecase.NewGetAddressBrasilAPIUseCase(*httpClient)
	processAddress := usecase.NewProcessAddressUseCase(ViaCep, BrasilAPI)

	address, err := processAddress.Execute(ctx, "09330340")
	if err != nil {
		log.Fatal(err)
	}
	jsonBytes, err := json.Marshal(address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonBytes))

}
