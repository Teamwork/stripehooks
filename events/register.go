package register

// registerEventHandlers is the place to add your handlers for different Stripe hooks
import (
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/event"
)

type StripeWebhook struct {
	Type     string
	Id       string
	LiveMode bool
	Data     stripe.EventData
}

type handleEvent func(event *stripe.Event) (err error)
type validateHook func(hook StripeWebhook) (event *stripe.Event, err error)

var EventHandlers = make(map[string]handleEvent)
var ValidationHandler validateHook

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
func RegisterEventHandler(event string, h handleEvent) {
	eventHandlers[event] = h
}

// Register your validation handler function here
// A typical validation is to call back to Stripe
// and verify and use the event that check returns
func RegisterValidationHandler(h validateHook) {
	validationHandler = h

	// If you want no validation with Stripe, then pass
	// validationHandler = func(hook StripeWebhook) (*stripe.Event, error) {
	//  return &stripe.Event{
	//     Data: &hook.Data,
	//     Live: hook.LiveMode,
	//     Type: hook.Type,
	//     ID:   hook.Id,
	// }, nil
	// }
}

// VerifyEventWithStripe verifies the event received via webhook with Stripe
// using the ID to confirm webhook is legit - extra security step
// Note: If the stripe webhook is not Livemode then this bypass the call to Stripe
// and use the event we have from the original webhook.
// A testmode hook will always fail this callback otherwise.
func VerifyEventWithStripe(hook StripeWebhook) (*stripe.Event, error) {
	if !hook.LiveMode {
		return &stripe.Event{
			Data: &hook.Data,
			Live: hook.LiveMode,
			Type: hook.Type,
			ID:   hook.Id,
		}, nil
	}
	stripe.Key = "" // your private stripe api key goes here
	return event.Get(hook.Id)
}
