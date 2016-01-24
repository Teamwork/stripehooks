package events

// registerEventHandlers is the place to add your handlers for different Stripe hooks
import (
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/event"
)

// StripeWebhook holds the data sent by Stripe through the webhook
type StripeWebhook struct {
	Type     string
	ID       string
	LiveMode bool
	Data     stripe.EventData
}

type eventHandler func(event *stripe.Event) (err error)
type validateHandler func(hook StripeWebhook) (event *stripe.Event, err error)

// EventHandlers is a mapping of stripe events to their handleEvent functions
var EventHandlers = make(map[string]handleEvent)

// ValidationHandler is an
var ValidationHandler validateHandler

// RegisterEventHandlers is the place to add your handlers for each
// event type you want to support from Stripe webhooks
func RegisterEventHandlers() {
	// Example handle incoming hook for invoice.payment_succeeded event
	registerEventHandler("invoice.payment_succeeded", func(event *stripe.Event) error {
		fmt.Printf("Event %s handled, type %s\n", event.ID, event.Type)
		return nil
	})

	registerEventHandler("invoice.payment_failed", func(event *stripe.Event) error {
		// Exmaple, send email to customer
		return nil
	})
}

// RegisterEventHandler adds the handler to the list of handlers
// It will overwrite handler if one exists for event already
func RegisterEventHandler(event string, handleFunc eventHandler) {
	eventHandlers[event] = handleFunc
}

// RegisterValidationHandler allows you to add your validation handler function here
// A typical validation is to call back to Stripe
// and verify and use the event that check returns.
//
// If you want no validation with Stripe, then just use the original event data e.g.
// validationHandler = func(hook StripeWebhook) (*stripe.Event, error) {
//      return &stripe.Event{
//          Data: &hook.Data,
//          Live: hook.LiveMode,
//          Type: hook.Type,
//          ID:   hook.Id,
//      }, nil
// }
func RegisterValidationHandler(handleFunc validateHandler) {
	validationHandler = handleFunc
}

// VerifyEventWithStripe verifies the event received via webhook with Stripe
// using the ID to confirm webhook is legit - extra security step
// Note: If the stripe webhook is not Livemode then this bypasses the call to Stripe
// and uses the event we have from the original webhook.
// Without that bypass, a testmode hook would always fail this callback
func VerifyEventWithStripe(hook StripeWebhook) (*stripe.Event, error) {
	if !hook.LiveMode {
		return &stripe.Event{
			Data: &hook.Data,
			Live: hook.LiveMode,
			Type: hook.Type,
			ID:   hook.ID,
		}, nil
	}
	stripe.Key = "" // your private stripe api key goes here
	return event.Get(hook.ID)
}
