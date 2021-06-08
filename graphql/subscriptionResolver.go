package graphql

import (
	"strings"

	graphqlws "github.com/functionalfoundry/graphqlws"
	graphql "github.com/graphql-go/graphql"
)

func subscriptionSubscriber(channel, payload string) {
	subscriptions := subscriptionManager.Subscriptions()
	for conn := range subscriptions {
		for _, subscription := range subscriptions[conn] {

			params := graphql.Params{
				Schema:         RustySchema, // The GraphQL schema
				RequestString:  subscription.Query,
				VariableValues: subscription.Variables,
				OperationName:  subscription.OperationName,
			}
			result := graphql.Do(params)
			data := graphqlws.DataMessagePayload{
				Data: result.Data,
			}
			if strings.Contains(subscription.Fields[0], "harmonyChange") {
				subscription.SendData(&data)
			} else if strings.Contains(subscription.Fields[0], "lutronChanges") {
				subscription.SendData(&data)
			} else if strings.Contains(subscription.Fields[0], "menuChange") {
				subscription.SendData(&data)
			} else if strings.Contains(subscription.Fields[0], "life360Change") {
				subscription.SendData(&data)
			}
		}
	}
}
