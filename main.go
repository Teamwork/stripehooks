package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go"
	"github.com/teamwork/stripehooks/events"
)

var (
	port int
	path string
)

func main() {

	flag.IntVar(&port, "port", 8080, "port number to bind the server to e.g. '8080' would be http://localhost:8080/<path>")
	flag.StringVar(&path, "path", "stripe", "the endpoint to recieve webhooks on e.g. 'stripe' would be http://localhost:<port>/stripe")
	flag.Parse()

	events.RegisterValidationHandler(events.VerifyEventWithStripe)
	events.RegisterEventHandlers()

	http.HandleFunc("/"+path, stripeHandler)

	fmt.Println("Listening...")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
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

	hook, err := parseBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	handler, ok := events.EventHandlers[hook.Type]
	if !ok {
		fmt.Fprintf(w, "No handler for this event %s", hook.Type)
		return
	}

	event, err := events.ValidationHandler(hook)
	if err != nil {
		http.Error(w, fmt.Sprintf("Verifying with stripe failed for event %s", hook.ID), http.StatusInternalServerError)
		return
	}

	err = handler(event)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Event handled successfully %s", event.ID)
	return
}

// parseBody unmarshals the payload from Stripe to our StripeWebhook struct
func parseBody(body io.Reader) (event *stripe.Event, err error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &event)
	return event, err
}
