package activitypub

import (
	"fmt"
	ap "github.com/go-ap/activitypub"
	as "github.com/go-ap/activitystreams"
	"github.com/go-ap/auth"
	"github.com/go-ap/errors"
	"github.com/go-ap/handlers"
	"github.com/pborman/uuid"
	"net/url"
	"path"
	"strings"
)

var ServiceIRI as.IRI

func Self(baseURL as.IRI) auth.Service {
	url, _ := baseURL.URL()
	inbox := *url
	inbox.Path = path.Join(inbox.Path, string(handlers.Inbox))
	outbox := *url
	outbox.Path = path.Join(outbox.Path, string(handlers.Outbox))

	oauth := *url
	oauth.Path = path.Join(oauth.Path, "oauth/")
	return auth.Service{
		Person: ap.Person{
			Parent: as.Person{
				ID:           as.ObjectID(url.String()),
				Type:         as.ServiceType,
				Name:         as.NaturalLanguageValues{{Ref: as.NilLangRef, Value: "self"}},
				AttributedTo: as.IRI("https://github.com/mariusor"),
				Audience:     as.ItemCollection{as.PublicNS},
				Content:      nil, //as.NaturalLanguageValues{{Ref: as.NilLangRef, Value: ""}},
				Icon:         nil,
				Image:        nil,
				Location:     nil,
				Summary:      as.NaturalLanguageValues{{Ref: as.NilLangRef, Value: "Generic ActivityPub service"}},
				Tag:          nil,
				URL:          baseURL,
			},
			Inbox:  as.IRI(inbox.String()),
			Outbox: as.IRI(outbox.String()),
			Endpoints: &ap.Endpoints{
				OauthAuthorizationEndpoint: as.IRI(fmt.Sprintf("%s/authorize", oauth.String())),
				OauthTokenEndpoint:         as.IRI(fmt.Sprintf("%s/token", oauth.String())),
			},
		},
	}
}

func DefaultServiceIRI(baseURL string) as.IRI {
	u, _ := url.Parse(baseURL)
	// TODO(marius): I don't like adding the / folder to something like http://fedbox.git
	// I need to find an
	if u.Path == "" {
		u.Path = "/"
	}
	return as.IRI(u.String())
}

// ItemByType
func ItemByType(typ as.ActivityVocabularyType) (as.Item, error) {
	if as.ActorTypes.Contains(typ) {
		return &auth.Person{Person: ap.Person{Parent: ap.Parent{Type: typ}}}, nil
	} else if as.ActivityTypes.Contains(typ) {
		return &as.Activity{Parent: as.Parent{Type: typ}}, nil
	} else if typ == as.CollectionType {
		return &Collection{Parent: as.Parent{Type: typ}}, nil
	} else if typ == as.OrderedCollectionType {
		return &OrderedCollection{Parent: as.Parent{Type: typ}}, nil
	}
	return ap.JSONGetItemByType(typ)
}

// ToOrderedCollection
func ToOrderedCollection(it as.Item) (*OrderedCollection, error) {
	switch o := it.(type) {
	case *OrderedCollection:
		return o, nil
	case OrderedCollection:
		return &o, nil
	case *as.OrderedCollection:
		col := OrderedCollection(*o)
		return &col, nil
	case as.OrderedCollection:
		col := OrderedCollection(o)
		return &col, nil
	}
	return nil, errors.Newf("invalid ordered collection")
}

// ToCollection
func ToCollection(it as.Item) (*Collection, error) {
	switch o := it.(type) {
	case *Collection:
		return o, nil
	case Collection:
		return &o, nil
	case *as.Collection:
		col := Collection(*o)
		return &col, nil
	case as.Collection:
		col := Collection(o)
		return &col, nil
	}
	return nil, errors.Newf("invalid  collection")
}

// GenerateID generates an unique identifier for the it ActivityPub Object.
func GenerateID(it as.Item, partOf string, by as.Item) (as.ObjectID, error) {
	uuid := uuid.New()
	id := as.ObjectID(fmt.Sprintf("%s/%s", strings.ToLower(partOf), uuid))
	if as.ActivityTypes.Contains(it.GetType()) {
		err := ap.OnActivity(it, func(a *as.Activity) error {
			a.ID = id
			return nil
		})
		if err != nil {
			return id, err
		}
	}
	if as.ActorTypes.Contains(it.GetType()) {
		err := auth.OnPerson(it, func(p *auth.Person) error {
			p.ID = id
			return nil
		})
		if err != nil {
			return id, err
		}
	}
	if as.ObjectTypes.Contains(it.GetType()) {
		switch it.GetType() {
		case as.PlaceType:
			p, err := as.ToPlace(it)
			if err != nil {
				return id, err
			}
			p.ID = id
			it = p
		case as.ProfileType:
			p, err := as.ToProfile(it)
			if err != nil {
				return id, err
			}
			p.ID = id
			it = p
		case as.RelationshipType:
			r, err := as.ToRelationship(it)
			if err != nil {
				return id, err
			}
			r.ID = id
			it = r
		case as.TombstoneType:
			t, err := as.ToTombstone(it)
			if err != nil {
				return id, err
			}
			t.ID = id
			it = t
		default:
			err := ap.OnObject(it, func(o *ap.Object) error {
				o.ID = id
				return nil
			})
			if err != nil {
				return id, err
			}
		}
	}
	return id, nil
}
