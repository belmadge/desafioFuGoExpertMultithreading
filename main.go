package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func main() {
	ctx := context.Background()
	ctxReq, cancelReq := context.WithCancel(ctx)

	viaCEPResponse := make(chan ViaCEPResponse)
	go requestViaCEP(ctxReq, viaCEPResponse)

	brasilAPIResponse := make(chan BrasilAPIResponse)
	go requestBrasilAPI(ctx, brasilAPIResponse)

	select {
	case res := <-viaCEPResponse:
		fmt.Printf("Via CEP Response: %s\n", res)
		ctxReq.Done()
	case res := <-brasilAPIResponse:
		fmt.Printf("Brasil API Response: %s\n", res)
		ctxReq.Done()
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: Timeout")
		cancelReq()
	}
}

func requestViaCEP(ctx context.Context, response chan<- ViaCEPResponse) {
	cep := "01153000"
	viaCEPURL := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", viaCEPURL, nil)
	if err != nil {
		fmt.Printf("Error in the request of ViaCEP: %s\n", err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error in the response of ViaCEP: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error in the body of ViaCEP: %s\n", err.Error())
		return
	}

	var bodyResponse ViaCEPResponse
	err = json.Unmarshal(body, &bodyResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling of ViaCEP %s\n", err.Error())
		return
	}

	response <- bodyResponse
}

func requestBrasilAPI(ctx context.Context, response chan<- BrasilAPIResponse) {
	cep := "01153000"
	brasilAPIURL := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	req, err := http.NewRequestWithContext(ctx, "GET", brasilAPIURL, nil)
	if err != nil {
		fmt.Printf("Error in the request of BrasilAPI: %s\n", err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error in the response of BrasilAPI: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error in the body of BrasilAPI: %s\n", err.Error())
		return
	}

	var bodyResponse BrasilAPIResponse
	err = json.Unmarshal(body, &bodyResponse)
	if err != nil {
		fmt.Printf("Error unmarshalling of BrasilAPI %s\n", err.Error())
		return
	}

	response <- bodyResponse
}
