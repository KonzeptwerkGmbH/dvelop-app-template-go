package domain

// BcConfig holds the connection settings for Business Central
type BcConfig struct {
	BaseURL  string
	Mandant  string
	Username string
	Password string
}

// BcConfigRepository persists BC connection settings
type BcConfigRepository interface {
	GetBcConfig() BcConfig
	SaveBcConfig(config BcConfig)
}

// Debitor represents a Business Central customer (Debitor)
type Debitor struct {
	No      string
	Name    string
	Address string
	City    string
	Phone   string
	Email   string
}

// DebitorClient fetches debtors from Business Central
type DebitorClient interface {
	GetDebitoren(config BcConfig) ([]Debitor, error)
}
