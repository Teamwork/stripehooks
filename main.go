package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/teamwork/stripehooks/events"
)

func main() {
	events.RegisterValidationHandler(events.VerifyEventWithStripe)
	events.RegisterEventHandlers()

	http.HandleFunc("/stripe", stripeHandler)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}

// stripeHandler is responsible for receiving the stripe webhook
// and then taking that request through verification and event handlers
// and return 200 OK to Stripe or some other HTTP response code in the event
// of an error. Note: Stripe will retry any events returned outside of 2XX.
func stripeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "route only supports POST", http.StatusBadRequest)
		return
	}

	d, err := parseBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handler, ok := events.EventHandlers[d.Type]
	if !ok {
		fmt.Fprintf(w, "No handler for this event %s", d.Type)
		return
	}

	event, err := events.ValidationHandler(d)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verifying with stripe failed for event %s", d.Id), http.StatusInternalServerError)
		return
	}

	err = handle(event)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Event handled successfully %s", d.ID)
	return
}

// parseBody unmarshals the payload from Stripe to our StripeWebhook struct
func parseBody(body io.Reader) (hook StripeWebhook, err error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &hook)
	return hook, err
}
