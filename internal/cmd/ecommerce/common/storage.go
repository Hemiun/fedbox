package common

import (
	vocab "github.com/go-ap/activitypub"
	st "github.com/go-ap/fedbox/storage"
	"github.com/go-ap/processing"
	"github.com/openshift/osin"
)

type Storage interface {
	ClientSaver
	ClientLister
	osin.Storage
	processing.Store
	st.PasswordChanger
}

type ClientSaver interface {
	// UpdateClient updates the client (identified by it's id) and replaces the values with the values of client.
	UpdateClient(c osin.Client) error
	// CreateClient stores the client in the database and returns an error, if something went wrong.
	CreateClient(c osin.Client) error
	// RemoveClient removes a client (identified by id) from the database. Returns an error if something went wrong.
	RemoveClient(id string) error
}

type ClientLister interface {
	// ListClients lists existing clients
	ListClients() ([]osin.Client, error)
	GetClient(id string) (osin.Client, error)
}

type MetadataTyper interface {
	LoadMetadata(vocab.IRI) (*processing.Metadata, error)
	SaveMetadata(processing.Metadata, vocab.IRI) error
}
