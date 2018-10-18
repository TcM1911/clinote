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
	cfg.UDB = db
	return evernote.NewClient(cfg)
}

func newClient(opts clinote.ClientOption) *clinote.Client {
	cfg := new(clinote.DefaultConfig)
	db, err := storage.Open(cfg.GetConfigFolder())
	if err != nil {
		panic("Error when opening the database: " + err.Error())
	}
	cfg.DB = db
	cfg.UDB = db
	ec := evernote.NewClient(cfg)
	ns, err := ec.GetNoteStore()
	if err != nil {
		panic("Error when getting notestore: " + err.Error())
	}
	return clinote.NewClient(cfg, db, ns, opts)
}
