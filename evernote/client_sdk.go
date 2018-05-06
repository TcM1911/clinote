package evernote

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/TcM1911/clinote/config"
	ec "github.com/TcM1911/evernote-sdk-golang/client"
	"github.com/TcM1911/evernote-sdk-golang/notestore"
	"github.com/mrjones/oauth"
)

var apiConsumer = "clinote"
var apiSecret = "e9a3234ceefed62b"
var setup sync.Once
var devBuild = false

// APIClient is the interface for the api client.
type APIClient interface {
	// GetNoteStore returns the note store for the user.
	GetNoteStore() (NotestoreClient, error)
	// GetAuthorizedToken gets the authorized token from the server.
	GetAuthorizedToken(tmpToken *oauth.RequestToken, verifier string) (token string, err error)
	// GetRequestToken requests a request token from the server.
	GetRequestToken(callbackURL string) (token *oauth.RequestToken, url string, err error)
	// GetConfig returns the client's configuration.
	GetConfig() config.Configuration
	// GetAPIToken returns the user's token.
	GetAPIToken() string
}

// Client is an implementation of the client interface for Evernote.
type Client struct {
	// Config holds all the configurations.
	Config config.Configuration
	// APIToken is the access token for the user's account.
	apiToken   string
	ns         NotestoreClient
	evernote   *ec.EvernoteClient
	evernoteNS *notestore.NoteStoreClient
}

// GetAPIToken is the access token for the user's account.
func (c Client) GetAPIToken() string {
	return c.apiToken
}

// GetConfig returns the configuration.
func (c *Client) GetConfig() config.Configuration {
	return c.Config
}

// GetNoteStore returns a notestore client for the user.
func (c *Client) GetNoteStore() (NotestoreClient, error) {
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
	store := &Notestore{client: c, evernoteNS: ns}
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
func NewClient(cfg config.Configuration) *Client {
	client := new(Client)
	client.Config = cfg
	setup.Do(func() {
		env := ec.PRODUCTION
		if devBuild {
			env = ec.SANDBOX
		}
		client.evernote = ec.NewClient(apiConsumer, apiSecret, env)
		devToken := os.Getenv("EVERNOTE_DEV_TOKEN")
		if devToken != "" {
			client.apiToken = devToken
		} else {
			cacheDir := cfg.GetCacheFolder()
			fp := filepath.Join(cacheDir, "session")
			if _, err := os.Stat(fp); os.IsNotExist(err) {
				return
			}
			f, err := os.OpenFile(fp, os.O_RDONLY, 0600)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err.Error())
			}
			client.apiToken = string(b)
		}
	})
	return client
}
