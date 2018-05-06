package client

import (
	"fmt"

	"github.com/TcM1911/evernote-sdk-golang/notestore"
	"github.com/TcM1911/evernote-sdk-golang/userstore"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/mrjones/oauth"
)

type EnvironmentType int

const (
	SANDBOX EnvironmentType = iota
	PRODUCTION
	YINXIANG
)

type EvernoteClient struct {
	host        string
	oauthClient *oauth.Consumer
	userStore   *userstore.UserStoreClient
}

func NewClient(key, secret string, envType EnvironmentType) *EvernoteClient {
	host := "www.evernote.com"
	if envType == SANDBOX {
		host = "sandbox.evernote.com"
	} else if envType == YINXIANG {
		host = "app.yinxiang.com"
	}
	client := oauth.NewConsumer(
		key, secret,
		oauth.ServiceProvider{
			RequestTokenUrl:   fmt.Sprintf("https://%s/oauth", host),
			AuthorizeTokenUrl: fmt.Sprintf("https://%s/OAuth.action", host),
			AccessTokenUrl:    fmt.Sprintf("https://%s/oauth", host),
		},
	)
	return &EvernoteClient{
		host:        host,
		oauthClient: client,
	}
}

func (c *EvernoteClient) GetRequestToken(callBackURL string) (*oauth.RequestToken, string, error) {
	return c.oauthClient.GetRequestTokenAndUrl(callBackURL)
}

func (c *EvernoteClient) GetAuthorizedToken(requestToken *oauth.RequestToken, oauthVerifier string) (*oauth.AccessToken, error) {
	return c.oauthClient.AuthorizeToken(requestToken, oauthVerifier)
}

func (c *EvernoteClient) GetUserStore() (*userstore.UserStoreClient, error) {
	if c.userStore != nil {
		return c.userStore, nil
	}
	evernoteUserStoreServerURL := fmt.Sprintf("https://%s/edam/user", c.host)
	evernoteUserTrans, err := thrift.NewTHttpPostClient(evernoteUserStoreServerURL)
	if err != nil {
		return nil, err
	}
	c.userStore = userstore.NewUserStoreClientFactory(
		evernoteUserTrans,
		thrift.NewTBinaryProtocolFactoryDefault(),
	)
	return c.userStore, nil
}

func (c *EvernoteClient) GetNoteStore(authenticationToken string) (*notestore.NoteStoreClient, error) {
	us, err := c.GetUserStore()
	if err != nil {
		return nil, err
	}
	notestoreURL, err := us.GetNoteStoreUrl(authenticationToken)
	if err != nil {
		return nil, err
	}
	ns, err := c.GetNoteStoreWithURL(notestoreURL)
	return ns, nil
}

func (c *EvernoteClient) GetNoteStoreWithURL(notestoreURL string) (*notestore.NoteStoreClient, error) {
	evernoteNoteTrans, err := thrift.NewTHttpPostClient(notestoreURL)
	if err != nil {
		return nil, err
	}
	client := notestore.NewNoteStoreClientFactory(
		evernoteNoteTrans,
		thrift.NewTBinaryProtocolFactoryDefault(),
	)
	return client, nil
}
