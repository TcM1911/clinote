package main

import (
	"github.com/TcM1911/clinote"
	"github.com/TcM1911/clinote/evernote"
	"github.com/TcM1911/clinote/storage"
)

func defaultClient() *evernote.Client {
	cfg := &clinote.DefaultConfig{}
	db, err := storage.Open(cfg.GetConfigFolder())
	if err != nil {
		panic("Error when opening the database: " + err.Error())
	}
	cfg.DB = db
	return evernote.NewClient(cfg)
}
