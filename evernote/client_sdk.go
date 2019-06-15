package evernote

import (
	"github.com/TcM1911/clinote"
	ec "github.com/TcM1911/evernote-sdk-golang/client"
	"github.com/TcM1911/evernote-sdk-golang/notestore"
	"github.com/mrjones/oauth"
)

var apiConsumer = "clinote"
var apiSecret = "e9a3234ceefed62b"
var devBuild = false

// Client is an implementation of the client interface for Evernote.
type Client struct {
	// Config holds all the configurations.
	Config clinote.Configuration
	// APIToken is the access token for the user's account.
	apiToken   string
	ns         clinote.NotestoreClient
	evernote   *ec.EvernoteClient
	evernoteNS *notestore.NoteStoreClient
}

// Close shuts down the client.
func (c *Client) Close() error {
	return c.Config.Close()
}

// GetAPIToken is the access token for the user's account.
func (c Client) GetAPIToken() string {
	return c.apiToken
}

// GetConfig returns the configuration.
func (c *Client) GetConfig() clinote.Configuration {
	return c.Config
}

// GetNoteStore returns a notestore client for the user.
func (c *Client) GetNoteStore() (clinote.NotestoreClient, error) {
	if c.ns != nil {
		return c.ns, nil
	}
	if c.apiToken == "" {
		return nil, ErrNotLoggedIn
	}
	ns, err := c.evernote.GetNoteStore(c.apiToken)
	if err != nil {
		return nil, err
	}
	c.evernoteNS = ns
	store := &Notestore{apiToken: c.apiToken, evernoteNS: ns}
	c.ns = store
	return store, nil
}

// GetAuthorizedToken gets the authorized token from the server.
func (c *Client) GetAuthorizedToken(tmpToken *oauth.RequestToken, verifier string) (string, error) {
	token, err := c.evernote.GetAuthorizedToken(tmpToken, verifier)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

// GetRequestToken requests a request token from the server.
func (c *Client) GetRequestToken(callback string) (*oauth.RequestToken, string, error) {
	return c.evernote.GetRequestToken(callback)
}

// NewClient creates a new Evernote client.
func NewClient(cfg clinote.Configuration) *Client {
	client := new(Client)
	client.Config = cfg
	env := ec.PRODUCTION

	key := migrateOldSession(cfg)
	if key != "" {
		// Migrate an old session.
		settings, err := cfg.Store().GetSettings()
		if err != nil {
			panic(err.Error())
		}
		settings.Credential = &clinote.Credential{
			Name:     "OAuth",
			Secret:   key,
			CredType: clinote.EvernoteCredential,
		}
		if err := cfg.Store().StoreSettings(settings); err != nil {
			panic(err.Error())
		}
	} else {
		settings, err := cfg.Store().GetSettings()
		if err != nil {
			panic(err.Error())
		}
		key = settings.APIKey
		if settings.Credential != nil && settings.Credential.CredType == clinote.EvernoteSandboxCredential {
			env = ec.SANDBOX
		}
	}
	client.evernote = ec.NewClient(apiConsumer, apiSecret, env)
	client.apiToken = key

	return client
}
