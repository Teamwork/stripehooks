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