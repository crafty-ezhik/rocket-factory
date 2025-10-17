package payment_test

import (
	"errors"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestPaymentSuccess() {
	var (
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		paymentMethod = "CARD"

		expectedUUID = uuid.New().String()
	)
	s.service.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).Return(expectedUUID, nil)

	transactionUUID, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().Equal(expectedUUID, transactionUUID)
}

func (s *ServiceSuite) TestPaymentFail() {
	var (
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		paymentMethod = "CARD"

		expectedUUID = ""
		ErrPaid      = errors.New("payment failed")
	)
	s.service.On("PayOrder", s.ctx, orderUUID, userUUID, paymentMethod).Return("", ErrPaid)

	transactionUUID, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, paymentMethod)
	s.Require().ErrorIs(err, ErrPaid)
	s.Require().Equal(expectedUUID, transactionUUID)
}
