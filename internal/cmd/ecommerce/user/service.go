package user

import (
	"git.sr.ht/~mariusor/lw"
	_ "git.sr.ht/~mariusor/lw"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	ap "github.com/go-ap/fedbox/activitypub"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-ap/filters"
	"time"
)

const (
// keyType    = fedbox.KeyTypeED25519
)

type UserService struct {
	db        common.Storage
	ctl       common.Control
	baseURL   string
	logger    lw.Logger
	appActors vocab.Item
}

func NewUserService(ctl common.Control, db common.Storage, baseURL string, l lw.Logger) (*UserService, error) {
	var target UserService
	var err error
	target.db = db
	target.ctl = ctl
	target.baseURL = baseURL
	target.logger = l

	baseIRI := vocab.IRI(baseURL)
	actorFilters := filters.FiltersNew()
	actorFilters.IRI = filters.ActorsType.IRI(baseIRI)
	actorFilters.Type = filters.CompStrs{filters.StringEquals("Application")}

	target.appActors, err = db.Load(actorFilters.GetLink())
	if err != nil {
		l.Errorf("can't init user service: %v", err)
		return nil, err
	}

	return &target, nil
}

func (s *UserService) AddUser(ur UserRequest, caller vocab.Actor) (vocab.Item, error) {
	var it vocab.Item

	if ur.Name == "" || ur.Password == "" {
		return it, errors.Errorf("User credentials doesn't pass")
	}
	if !s.isSuperUser(caller) {
		err := errors.Errorf("Actor has insufficient privileges")
		return it, errors.NewForbidden(err, "Access denied")
	}

	authIRI := caller.GetLink()
	author, err := ap.LoadActor(s.db, caller.GetLink())
	if err != nil {
		s.logger.Errorf("Can't load author from db", err)
		return it, err
	}

	tags := make(vocab.ItemCollection, 0)

	objectsCollection := filters.ObjectsType.IRI(vocab.IRI(s.baseURL))
	allObjects, _ := s.db.Load(objectsCollection)

	vocab.OnCollectionIntf(allObjects, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			vocab.OnObject(it, func(object *vocab.Object) error {
				for _, tag := range ur.Tags {
					if object.Name.First().Value.String() != tag {
						continue
					}
					if object.AttributedTo.GetLink() != authIRI {
						continue
					}
					tags.Append(object)
				}
				return nil
			})
		}
		return nil
	})

	typ := vocab.PersonType

	now := time.Now().UTC()
	newPerson := &vocab.Person{
		Type:         typ,
		AttributedTo: author.GetLink(),
		Generator:    author.GetLink(),
		Published:    now,
		Summary: vocab.NaturalLanguageValues{
			{vocab.NilLangRef, vocab.Content(ur.Comments)},
		},
		Updated: now,
		PreferredUsername: vocab.NaturalLanguageValues{
			{vocab.NilLangRef, vocab.Content(ur.Name)},
		},
	}
	if len(tags) > 0 {
		newPerson.Tag = tags
	}

	if newPerson, err = s.ctl.AddActor(newPerson, []byte(ur.Password), &author); err != nil {
		s.logger.Errorf("Can't save new caller", err)
		return it, err
	}

	//fmt.Printf("Added %q [%s]: %s\n", typ, name, newPerson.GetLink())
	//
	//	if metaSaver, ok := s.ctl.Storage.(s.MetadataTyper); ok {
	//		if err := AddKeyToItem(metaSaver, newPerson, keyType); err != nil {
	//			Errf("Error saving metadata for %s: %s", name, err)
	//		}
	//	}
	//
	return newPerson, nil
}

func (s *UserService) DeleteUser(caller vocab.Actor, actorID string) error {
	// check access for the caller actor
	if !s.isSuperUser(caller) {
		err := errors.Errorf("Actor has insufficient privileges")
		return errors.NewForbidden(err, "Access denied")
	}

	// Prepare actorIRI
	f := filters.FiltersNew()
	f.IRI = filters.ActorsType.IRI(vocab.IRI(s.baseURL)).AddPath(actorID)

	// Load actor from repo
	actor, err := ap.LoadActor(s.db, f.GetLink())

	// It's important to check that actor is not empty. If your pass empty entity into db.Delete it will delete entire storage!
	if err != nil || !vocab.IsNotEmpty(actor) {
		s.logger.Errorf("Can't load author from db", err)
		return errors.NewNotFound(err, "actor not found")
	}

	err = s.db.Delete(actor)
	if err != nil {
		s.logger.Errorf("Can't delete author from db", err)
		return err
	}
	return nil
}

func (s *UserService) isSuperUser(actor vocab.Actor) bool {
	// TODO:
	var flIsSuperUser bool
	vocab.OnCollectionIntf(s.appActors, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			vocab.OnObject(it, func(act *vocab.Object) error {
				if act.ID == actor.AttributedTo {
					flIsSuperUser = true
				}
				return nil
			})
		}
		return nil
	})
	return flIsSuperUser
}
