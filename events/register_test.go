package events

import (
	"testing"

	"github.com/stripe/stripe-go"
)

func TestVerifyEventWithStripeWithTestMode(t *testing.T) {
	e := &stripe.Event{
		Type: "invoice",
		ID:   "in_123456789",
		Live: false,
		Data: &stripe.EventData{},
	}

	event, err := VerifyEventWithStripe(e)
	if err != nil {
		t.Fail()
	}
	if event.ID != e.ID && event.Type != e.Type {
		t.Errorf("Expected eventID to be %s but was %s", e.ID, event.ID)
	}
}
