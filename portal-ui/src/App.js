import React, { useLayoutEffect, useState } from "react";
import {
  ApolloClient,
  InMemoryCache,
  ApolloProvider,
  HttpLink,
  split,
  ApolloLink,
} from "@apollo/client";
import { onError } from "@apollo/client/link/error";
import { RetryLink } from "@apollo/client/link/retry";
import { getMainDefinition } from "@apollo/client/utilities";
import { WebSocketLink } from "@apollo/link-ws";
import './App.css';
import { render } from "react-dom";
import { AreaActivity } from "./pages/AreaActivity/AreaActivity";
import { makeStyles } from '@material-ui/core/styles';
import Icon from "@mdi/react";
import SideNav, { NavItem, NavIcon, NavText } from '@trendmicro/react-sidenav';
import './index.css';
import { mdiCog, mdiHome } from '@mdi/js';
import ClickOutside from "./componentLibrary/ClickOutside/click-outside";
import { BrowserRouter as Router, Route } from "react-router-dom";
import { Integration } from "./pages/Integration/Integration";

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

const linkToRetry = new RetryLink({
  delay: {
    initial: 300,
    max: Infinity,
    jitter: true
  },
  attempts: {
    max: 5,
    retryIf: (error, _operation) => !!error
  }
});

const client = new ApolloClient({
  cache: new InMemoryCache(),
  link: ApolloLink.from([
    onError((error) => {
      console.log({ error })
    }),
    linkToRetry,
    link
  ]),
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
  const [width] = useWindowSize();
  const [expanded, changeExpanded] = useState(false);

  return (
    <div className={width > 1100 ? classes.root : classes.tablet}>
      <ApolloProvider client={client}>
        {width > 1100 ? (
          <Router>
            <Route render={({ location, history }) => (
              <React.Fragment>
                <ClickOutside onClickOutside={() => {
                  changeExpanded(false);
                }}>
                  <SideNav
                    expanded={expanded}
                    onToggle={(expanded) => {
                      changeExpanded(expanded);
                    }}
                    onSelect={(selected) => {
                      const to = '/' + selected;
                      if (location.pathname !== to) {
                        history.push(to);
                      }
                    }}
                  >
                    {/* <SideNav.Toggle /> */}
                    <SideNav.Nav defaultSelected="home">
                      <NavItem eventKey="home">
                        <NavIcon>
                          <Icon path={mdiHome} size={2} />
                        </NavIcon>
                        <NavText>
                          Home
                        </NavText>
                      </NavItem>
                      <NavItem eventKey="settings">
                        <NavIcon>
                          <Icon path={mdiCog} size={2} />
                        </NavIcon>
                        <NavText>
                          Settings
                        </NavText>
                        <NavItem eventKey="settings/integrations">
                          <NavText>
                            Integrations
                          </NavText>
                        </NavItem>
                        <NavItem eventKey="settings/about">
                          <NavText>
                            About
                          </NavText>
                        </NavItem>
                      </NavItem>
                    </SideNav.Nav>
                  </SideNav>
                </ClickOutside>
                <main>
                  <Route path="/" exact component={() => <AreaActivity></AreaActivity>} />
                  <Route path="/home" component={() => <AreaActivity></AreaActivity>} />
                  <Route path="/settings" component={props => <div />} />
                  <Route path="/settings/integrations" component={props => <Integration />} />
                  <Route path="/settings/about" component={props => <div />} />
                </main>
              </React.Fragment>
            )} />
          </Router>) : (<div></div>)
        }
      </ApolloProvider>
    </div>
  );
}

render(<App />, document.getElementById("root"));

export default App;


// React Router v4 with React v16
// https://reactrouter.com/web/api/Route

{/* <Router>
    <Route render={({ location, history }) => (
        <React.Fragment>
            <SideNav
                onSelect={(selected) => {
                    const to = '/' + selected;
                    if (location.pathname !== to) {
                        history.push(to);
                    }
                }}
            >
                <SideNav.Toggle />
                <SideNav.Nav defaultSelected="home">
                    <NavItem eventKey="home">
                        <NavIcon>
                            <i className="fa fa-fw fa-home" style={{ fontSize: '1.75em' }} />
                        </NavIcon>
                        <NavText>
                            Home
                        </NavText>
                    </NavItem>
                    <NavItem eventKey="devices">
                        <NavIcon>
                            <i className="fa fa-fw fa-device" style={{ fontSize: '1.75em' }} />
                        </NavIcon>
                        <NavText>
                            Devices
                        </NavText>
                    </NavItem>
                </SideNav.Nav>
            </SideNav>
            <main>
                <Route path="/" exact component={props => <RootComponent />} />
                <Route path="/home" component={props => <Home />} />
                <Route path="/devices" component={props => <Devices />} />
            </main>
        </React.Fragment>
    )}
    />
</Router> */}