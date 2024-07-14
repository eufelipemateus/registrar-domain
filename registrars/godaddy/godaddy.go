package registrars

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const URL_DEV = "https://api.ote-godaddy.com"
const URL_PROD = "https://api.godaddy.com"

type APIKey struct {
	Key    string
	Secret string
}

type Registrar struct {
	Domain       string
	Period       int
	ApiKey       APIKey
	IsProduction bool `json:"isProduction"`
}

// RenewRequest é a estrutura para o corpo da requisição de renovação
type renewRequest struct {
	Period int `json:"period"`
}

// RenewResponse é a estrutura para a resposta da API
type RenewResponse struct {
	OrderID int64  `json:"orderId"`
	Status  string `json:"status"`
}

func (r Registrar) Renew() error {
	var url string

	if r.IsProduction {
		url = fmt.Sprintf("%s/v1/domains/%s/renew", URL_PROD, r.Domain)
	} else {
		url = fmt.Sprintf("%s/v1/domains/%s/renew", URL_DEV, r.Domain)
	}
	renewReq := renewRequest{Period: r.Period}
	reqBody, err := json.Marshal(renewReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Criar uma nova requisição HTTP POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Adicionar os cabeçalhos necessários
	req.Header.Set("Authorization",fmt.Sprintf( "sso-key %s:%s", r.ApiKey.Key, r.ApiKey.Secret))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Ler e processar a resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Verificar se a requisição foi bem-sucedida
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Deserializar a resposta
	var renewResp RenewResponse
	if err := json.Unmarshal(body, &renewResp); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return nil
}
