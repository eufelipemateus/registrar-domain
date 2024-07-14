package name

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const URL_DEV = "api.dev.name.com"
const URL_PROD = "api.name.com"

type APIKey struct {
	Key    string
	Secret string
}

type Registrar struct {
	Domain       string
	Period       int
	APIKey       APIKey
	IsProduction bool
}

// RenewRequest é a estrutura para o corpo da requisição de renovação
type renewRequest struct {
	Years int `json:"years"`
}

// RenewResponse é a estrutura para a resposta da API
type RenewResponse struct {
	OrderID   string `json:"orderId"`
	Domain    string `json:"domain"`
	ExpiresAt string `json:"expiresAt"`
}

func (r Registrar) Renew() error {
	var url string

	if r.IsProduction {
		url = fmt.Sprintf("https://%s:%s@%s/v4/domains/%s:renew", r.APIKey.Key, r.APIKey.Secret, URL_PROD, r.Domain)
	} else {
		url = fmt.Sprintf("https://%s:%s@%s/v4/domains/%s:renew", r.APIKey.Key, r.APIKey.Secret, URL_DEV, r.Domain)
	}

	// Criar o corpo da requisição
	renewReq := renewRequest{Years: r.Period}
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
	req.Header.Set("Content-Type", "application/json")

	// Criar um cliente HTTP e definir um timeout
	client := &http.Client{Timeout: 10 * time.Second}

	// Enviar a requisição
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

	fmt.Printf("Domain %s renewed successfully. Order ID: %s, Expires At: %s\n", renewResp.Domain, renewResp.OrderID, renewResp.ExpiresAt)
	return nil
}
