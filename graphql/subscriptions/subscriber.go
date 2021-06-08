package subscriptions

type Subscriber interface {
	AlertSubscription(data interface{}, subscriptionType string)
}
