package payment_gateway

type PaymentGateway interface {
	CreatePayment(orderID string, amount float64, method string) (*PaymentResult, error)
	VerifyCallback(payload []byte, signature string) (*CallbackPayload, error)
}

type PaymentResult struct {
	PaymentURL string
	GatewayRef string
}

type CallbackPayload struct {
	ExternalID    string
	Status        string
	TransactionID string
}

type Config struct {
	ServerKey  string
	ClientKey  string
	BaseURL    string
	Production bool
}
