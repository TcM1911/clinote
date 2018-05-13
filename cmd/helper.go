package cmd

import (
	"github.com/TcM1911/clinote/config"
	"github.com/TcM1911/clinote/evernote"
)

func defaultClient() *evernote.Client {
	cfg := &config.DefaultConfig{}
	return evernote.NewClient(cfg)
}
