import React, { useLayoutEffect, useState } from "react";
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
import Icon from "@mdi/react";
import SideNav, { Toggle, Nav, NavItem, NavIcon, NavText } from '@trendmicro/react-sidenav';
import './index.css';
import { mdiCog, mdiHome, mdiLightbulb } from '@mdi/js';
import ClickOutside from "./componentLibrary/ClickOutside/click-outside";

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: "#02464d",
    width: "100%",
    height: "100%",
    paddingLeft: '2rem'
    // position: "absolute",
  },
  tablet: {
    backgroundColor: "#02464d",
    width: "100%",
    height: "100%"
    // position: "absolute",
  },
  important: {
  }
}));

const rustyServer = "http://192.168.1.19:30919/graphql";
const rustyWs = "ws://192.168.1.19:30919/subscriptions";
// const localServer = "http://localhost:4000/graphql";
// const localWs = "ws://localhost:4000/subscriptions";

const httpLink = new HttpLink({
  uri: rustyServer,
});

// Create a WebSocket link:
const wsLink = new WebSocketLink({
  uri: rustyWs,
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

const defaultOptions = {
  watchQuery: {
    fetchPolicy: 'no-cache',
    errorPolicy: 'ignore',
  },
  query: {
    fetchPolicy: 'no-cache',
    errorPolicy: 'all',
  },
}

const client = new ApolloClient({
  cache: new InMemoryCache(),
  link,
  defaultOptions: defaultOptions,
});

function useWindowSize() {
  const [size, setSize] = useState([0, 0]);
  useLayoutEffect(() => {
    function updateSize() {
      setSize([window.innerWidth, window.innerHeight]);
    }
    window.addEventListener('resize', updateSize);
    updateSize();
    return () => window.removeEventListener('resize', updateSize);
  }, []);
  return size;
}

const App = () => {
  const classes = useStyles();
  const [width, height] = useWindowSize();
  const [expanded, changeExpanded] = useState(false);

  console.log({width})
  return (
    <div className={width > 1100 ? classes.root : classes.tablet}>
      <ApolloProvider client={client}>
        { width > 1100 ? (
          <ClickOutside onClickOutside={() => {
            changeExpanded(false);
          }}>
            <SideNav
              expanded={expanded}
              onToggle={(expanded) => {
                changeExpanded(expanded);
              }}
              onSelect={(selected) => {
                // Add your code here
              }}
            >
              <SideNav.Toggle />
              <SideNav.Nav defaultSelected="home">
                <NavItem eventKey="home">
                  <NavIcon>
                    <Icon path={mdiHome} size={2} />
                  </NavIcon>
                  <NavText>
                    Home
                  </NavText>
                </NavItem>
                <NavItem eventKey="charts">
                  <NavIcon>
                    <Icon path={mdiCog} size={2} />
                  </NavIcon>
                  <NavText>
                    Settings
                  </NavText>
                  <NavItem eventKey="charts/linechart">
                    <NavText>
                      Integrations
                    </NavText>
                  </NavItem>
                  <NavItem eventKey="charts/barchart">
                    <NavText>
                      Bar Chart
                    </NavText>
                  </NavItem>
                </NavItem>
              </SideNav.Nav>
            </SideNav>
          </ClickOutside>) : (<div></div>)


        }

        <AreaActivity></AreaActivity>
      </ApolloProvider>
    </div>
  );
}

render(<App />, document.getElementById("root"));

export default App;
