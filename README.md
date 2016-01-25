[![Build Status](https://travis-ci.org/Teamwork/stripehook.svg?branch=master)](https://travis-ci.org/Teamwork/stripehook)

# stripehook

A server to handle [Stripe webhooks](https://stripe.com/docs/webhooks). Add your own custom event handlers. 

## Run
`go get github.com/teamwork/stripehook`

`go run main.go` or `go build` && `./stripehook`

It accepts a port and a path part for configuring the endpoint e.g. `./stripehook --help`

```
Usage of ./stripehook:
  -path="stripe": the endpoint to recieve webhooks on e.g. 'stripe' would be http://localhost:<port>/stripe
  -port=8080: port number to bind the server to e.g. '8080' would be http://localhost:8080/<path>
```

## Example

In `events/register.go` add your event handlers to `RegisterEventHandlers()` function for any of the various events that Stripe send to your endpoint. You can see the list of events Stripe sends [here](https://stripe.com/docs/api#event_types). Here is an example of handling an `invoice.payment_succeeded` event. 

```
    // Example handle incoming hook for invoice.payment_succeeded event
    registerEventHandler("invoice.payment_succeeded", func(event *stripe.Event) error {
        fmt.Printf("Event %s handled, type %s\n", event.ID, event.Type)
        return nil
    })
```

`eventHandlers` is a `map[string]eventHandler` that holds all the events you register. There is also a `validationHandler` for checking the event Stripe sends down. This is typically either checking the event back with Stripe (extra step for security) or to use the event as-is and trust it. Examples of both are provided

```go
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
```


```go
validationHandler = func(hook StripeWebhook) (*stripe.Event, error) {
     return &stripe.Event{
         Data: &hook.Data,
         Live: hook.LiveMode,
         Type: hook.Type,
         ID:   hook.Id,
     }, nil
}
```