package events

// registerEventHandlers is the place to add your handlers for different Stripe hooks
import (
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/event"
)

type eventHandler func(hook *stripe.Event) (err error)
type validateHandler func(hook *stripe.Event) (event *stripe.Event, err error)

// EventHandlers is a mapping of stripe events to their handleEvent functions
var EventHandlers = make(map[string]eventHandler)

// ValidationHandler is an
var ValidationHandler validateHandler

// RegisterEventHandlers is the place to add your handlers for each
// event type you want to support from Stripe webhooks
func RegisterEventHandlers() {
	// Example handle incoming hook for invoice.payment_succeeded event
	RegisterEventHandler("invoice.payment_succeeded", func(event *stripe.Event) error {
		fmt.Printf("Event %s handled, type %s\n", event.ID, event.Type)
		return nil
	})

	RegisterEventHandler("invoice.payment_failed", func(event *stripe.Event) error {
		// Exmaple, send email to customer
		return nil
	})
}

// RegisterEventHandler adds the handler to the list of handlers
// It will overwrite handler if one exists for event already
func RegisterEventHandler(event string, handleFunc eventHandler) {
	EventHandlers[event] = handleFunc
}

// RegisterValidationHandler allows you to add your validation handler function here
// A typical validation is to call back to Stripe
// and verify and use the event that check returns.
//
// If you want no validation with Stripe, then just use the original event data e.g.
// validationHandler = func(hook *stripe.Event) (*stripe.Event, error) {
//      return &stripe.Event{
//          Data: &hook.Data,
//          Live: hook.LiveMode,
//          Type: hook.Type,
//          ID:   hook.Id,
//      }, nil
// }
func RegisterValidationHandler(handleFunc validateHandler) {
	ValidationHandler = handleFunc
}

// VerifyEventWithStripe verifies the event received via webhook with Stripe
// using the ID to confirm webhook is legit - extra security step
// Note: If the stripe webhook is not Livemode then this bypasses the call to Stripe
// and uses the event we have from the original webhook.
// Without that bypass, a testmode hook would always fail this callback
func VerifyEventWithStripe(hook *stripe.Event) (*stripe.Event, error) {
	if !hook.Live {
		return hook, nil
	}
	stripe.Key = "" // your private stripe api key goes here
	return event.Get(hook.ID)
}
