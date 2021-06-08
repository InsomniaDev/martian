package graphql

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/friendsofgo/graphiql"
	graphqlws "github.com/functionalfoundry/graphqlws"
	graphql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	integ "github.com/insomniadev/martian/integrations"
	"github.com/insomniadev/martian/modules/redispub"
	"github.com/rs/cors"
)

var (
	//Integrations is the main way for working with the smart home integrations
	Integrations        integ.Integrations
	subscriptionManager graphqlws.SubscriptionManager
	RustySchema         graphql.Schema
)

func initModels() {
	Integrations.Init()
}

type reqBody struct {
	Query string `json:"query"`
}

// Graphql is the main graphql entry point
func Graphql() {

	initModels()
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/graphql")
	if err != nil {
		panic(err)
	}
	RustySchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:        rootQuery,
		Mutation:     rootMutation,
		Subscription: rootSubscription,
	})
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:     &RustySchema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	corsHandler := cors.Default().Handler(h)

	subscriptionManager = graphqlws.NewSubscriptionManager(&RustySchema)

	graphqlwsHandler := graphqlws.NewHandler(graphqlws.HandlerConfig{
		SubscriptionManager: subscriptionManager,
	})

	redispub.NewSubscriber("subscriptions", subscriptionSubscriber)
	// The handler integrates seamlessly with existing HTTP servers
	http.Handle("/subscriptions", graphqlwsHandler)
	http.Handle("/graphql", corsHandler)
	http.Handle("/graphiql", graphiqlHandler)
	http.HandleFunc("/zwave", zwavehandler)
	http.ListenAndServe(":4000", nil)
}

func zwavehandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
