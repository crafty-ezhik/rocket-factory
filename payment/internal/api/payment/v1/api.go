package v1

import (
	"github.com/crafty-ezhik/rocket-factory/payment/internal/service"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

type API struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *API {
	return &API{paymentService: paymentService}
}
