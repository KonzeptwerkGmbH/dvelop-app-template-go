package bc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/d-velop/dvelop-app-template-go/domain"
)

type client struct {
	http *http.Client
}

func NewClient() domain.DebitorClient {
	return &client{
		http: &http.Client{Timeout: 15 * time.Second},
	}
}

type odataResponse struct {
	Value []bcCustomer `json:"value"`
}

type bcCustomer struct {
	No      string `json:"No"`
	Name    string `json:"Name"`
	Address string `json:"Address"`
	City    string `json:"City"`
	PhoneNo string `json:"Phone_No"`
	EMail   string `json:"E_Mail"`
}

func (c *client) GetDebitoren(config domain.BcConfig) ([]domain.Debitor, error) {
	endpoint := fmt.Sprintf("%s/ODataV4/Company('%s')/Customer",
		config.BaseURL,
		url.PathEscape(config.Mandant),
	)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Erstellen der Anfrage: %w", err)
	}
	req.SetBasicAuth(config.Username, config.Password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("verbindung zu Business Central fehlgeschlagen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Business Central antwortete mit Status %d", resp.StatusCode)
	}

	var odataResp odataResponse
	if err := json.NewDecoder(resp.Body).Decode(&odataResp); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen der Antwort: %w", err)
	}

	debitoren := make([]domain.Debitor, len(odataResp.Value))
	for i, cust := range odataResp.Value {
		debitoren[i] = domain.Debitor{
			No:      cust.No,
			Name:    cust.Name,
			Address: cust.Address,
			City:    cust.City,
			Phone:   cust.PhoneNo,
			Email:   cust.EMail,
		}
	}
	return debitoren, nil
}
