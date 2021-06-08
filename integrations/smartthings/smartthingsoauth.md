package smartthings

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
)

// NewOAuthConfig creates a new oauth2.config structure with the
// correct parameters to use smartthings.
func (st *SmartThings) NewOAuthConfig(client, secret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     client,
		ClientSecret: secret,
		Scopes:       []string{"app"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://graph.api.smartthings.com/oauth/authorize",
			TokenURL: "https://graph.api.smartthings.com/oauth/token",
		},
	}
}

// NewAuth creates a new Auth struct
func NewAuth(port int, config *oauth2.Config) (*Auth, error) {
	rnd, err := randomString(16)
	if err != nil {
		return nil, err
	}

	return &Auth{
		port:             port,
		config:           config,
		rchan:            make(chan oauthReturn),
		oauthStateString: rnd,
	}, nil
}

// FetchOAuthToken sets up the handler and a local HTTP server and fetches an
// Oauth token from the smartthings website.
func (g *Auth) FetchOAuthToken() (*oauth2.Token, error) {
	http.HandleFunc(rootPath, g.handleMain)
	http.HandleFunc(donePath, g.handleDone)
	http.HandleFunc(callbackPath, g.handleOAuthCallback)

	go http.ListenAndServe(":"+strconv.Itoa(g.port), nil)

	// Block on the return channel (this is set by handleOauthCallback)
	ret := <-g.rchan
	return ret.token, ret.err
}

// handleMain redirects the user to the main authentication page.
func (g *Auth) handleMain(w http.ResponseWriter, r *http.Request) {
	url := g.config.AuthCodeURL(g.oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleError shows a page indicating the authentication has failed.
func (g *Auth) handleError(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, authError)
}

// handleDone shows a page indicating the authentication is finished.
func (g *Auth) handleDone(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, authDone)
}

// handleOauthCallback fetches the callback from the OAuth provider and parses
// the URL, extracting the code and then exchanging it for a token.
func (g *Auth) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	// Make sure we have the same "state" as our request.
	state := r.FormValue("state")
	if state != g.oauthStateString {
		g.rchan <- oauthReturn{
			token: nil,
			err:   fmt.Errorf("invalid oauth state, expected %q, got %q", g.oauthStateString, state),
		}
		return
	}

	// Retrieve the code from the URL, and exchange for a token
	code := r.FormValue("code")
	token, err := g.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		g.rchan <- oauthReturn{
			token: nil,
			err:   fmt.Errorf("code exchange failed: %q", err),
		}
		return
	}

	// Return token.
	g.rchan <- oauthReturn{
		token: token,
		err:   nil,
	}
	// Redirect user to "Authentication done" page
	http.Redirect(w, r, donePath, http.StatusTemporaryRedirect)
	return
}

// GetEndPointsURI returns the smartthing endpoints URI. The endpoints
// URI is the base for all app requests.
func (st *SmartThings) GetEndPointsURI() (string, error) {
	// Fetch the JSON containing our endpoint URI
	resp, err := st.Client.Get(endPointsURI)
	if err != nil {
		return "", fmt.Errorf("error getting endpoints URI: %q", err)
	}
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading endpoints URI data: %q", err)
	}
	resp.Body.Close()
	if string(contents) == "[]" {
		return "", fmt.Errorf("endpoint URI returned no content")
	}

	// Only URI is fetched from JSON string.
	var ep []endpoints
	err = json.Unmarshal(contents, &ep)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON: %q", err)
	}
	return ep[0].URI, nil
}

// randomString generates a random string of bytes of the specified size
// and returns the its hexascii representation.
func randomString(size int) (string, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// GetToken returns the token for the ClientID and Secret specified in config.
// The function attempts to load the token from tokenFile first, and failing
// that, starts a full token authentication cycle with SmartThings. If
// tokenFile is blank, the function uses a default name under the current
// user's home directory. The token is saved to local disk before being
// returned to the caller.
//
// This function represents the most common (and possibly convenient) way to
// retrieve a token for a given ClientID and Secret.
func (st *SmartThings) GetToken(config *oauth2.Config) (*oauth2.Token, error) {
	// Attempt to load token from local storage. Fallback to full auth cycle.
	token, err := st.LoadToken()
	if err != nil || !token.Valid() {
		if config.ClientID == "" || config.ClientSecret == "" {
			return nil, errors.New("need ClientID and secret to generate a new token")
		}
		gst, err := NewAuth(defaultPort, config)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Please login by visiting http://localhost:%d\n", defaultPort)
		token, err = gst.FetchOAuthToken()
		if err != nil {
			return nil, err
		}

		// Once we have the token, save it locally for future use.
		err = st.SaveToken(token)
		if err != nil {
			return nil, err
		}
	}
	return token, nil
}
