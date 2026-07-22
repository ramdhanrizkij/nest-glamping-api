package payment_gateway

import (
	"fmt"
)

type ManualGateway struct{}

func NewManualGateway() PaymentGateway {
	return &ManualGateway{}
}

func (g *ManualGateway) CreatePayment(orderID string, amount float64, method string) (*PaymentResult, error) {
	return &PaymentResult{
		PaymentURL: fmt.Sprintf("/payments/%s/instructions", orderID),
		GatewayRef: fmt.Sprintf("MANUAL-%s", orderID),
	}, nil
}

func (g *ManualGateway) VerifyCallback(payload []byte, signature string) (*CallbackPayload, error) {
	return &CallbackPayload{
		Status: "paid",
	}, nil
}
