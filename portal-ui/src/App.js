import React from "react";
import {
  ApolloClient,
  InMemoryCache,
  ApolloProvider,
  HttpLink,
  split,
} from "@apollo/client";
import { getMainDefinition } from "@apollo/client/utilities";
import { WebSocketLink } from "@apollo/link-ws";
import './App.css';
import { render } from "react-dom";
import { AreaActivity } from "./pages/AreaActivity/AreaActivity";
import { makeStyles } from '@material-ui/core/styles';

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: "#02464d",
    width: "100%", 
    height: "100%",
    // position: "absolute",
  },
}));

const rustyServer = "http://192.168.1.19:30919/graphql";
const rustyWs = "ws://192.168.1.19:30919/subscriptions";
const localServer = "http://localhost:4000/graphql";
const localWs = "ws://localhost:4000/subscriptions";

const httpLink = new HttpLink({
  uri: localServer,
});

// Create a WebSocket link:
const wsLink = new WebSocketLink({
  uri: localWs,
  options: {
    reconnect: true,
  },
});

const link = split(
  // split based on operation type
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === "OperationDefinition" &&
      definition.operation === "subscription"
    );
  },
  wsLink,
  httpLink
);

const client = new ApolloClient({
  cache: new InMemoryCache(),
  link,
});

const App = () => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <ApolloProvider client={client}>
        <AreaActivity></AreaActivity>
      </ApolloProvider>
    </div>
  );
}

render(<App />, document.getElementById("root"));

export default App;
