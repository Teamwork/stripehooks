[![Build Status](https://travis-ci.org/Teamwork/stripehook.svg?branch=master)](https://travis-ci.org/Teamwork/stripehook)

# stripehook

A server to handle [Stripe webhooks](https://stripe.com/docs/webhooks). Add your own custom event handlers. 

## Run
`go get github.com/teamwork/stripehook`

`cd stripehook`

`go run main.go`

There is an upstart script provided to start/stop the server on Ubuntu, and have it recover from a crash automatically.  

## Example

In `events/register.go` add your event handlers to `RegisterEventHandlers()` function for any of the various events that Stripe send to your endpoint. You can see the list of events Stripe sends [here](https://stripe.com/docs/api#event_types). Here is an example of handling an `invoice.payment_succeeded` event. 

```
    // Example handle incoming hook for invoice.payment_succeeded event
    registerEventHandler("invoice.payment_succeeded", func(event *stripe.Event) error {
        fmt.Printf("Event %s handled, type %s\n", event.ID, event.Type)
        return nil
    })
```