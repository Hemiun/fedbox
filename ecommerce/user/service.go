package user

import (
	"git.sr.ht/~mariusor/lw"
	_ "git.sr.ht/~mariusor/lw"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	"github.com/go-ap/fedbox"
	ap "github.com/go-ap/fedbox/activitypub"
	"github.com/go-ap/fedbox/internal/cmd"
	"github.com/go-ap/filters"
	"time"
)

type UserService struct {
	ctl     *cmd.Control
	baseURL string
	logger  lw.Logger
}

func NewUserService(ctl *cmd.Control, baseURL string, l lw.Logger) *UserService {
	var target UserService
	target.ctl = ctl
	target.baseURL = baseURL
	target.logger = l
	return &target
}

const (
	keyType = fedbox.KeyTypeED25519
)

func (s *UserService) clientUri() vocab.IRI {
	// get URI for main client
	//TODO:
	return vocab.IRI("")
}

func (s *UserService) NewUser(ur UserRequest) (vocab.IRI, error) {
	if ur.Name == "" || ur.Password == "" {
		return "", errors.Errorf("User credentials doesn't pass")
	}

	authIRI := s.clientUri()

	if authIRI == "" {
		s.logger.Errorf("Can't get client actor Uri")
		return "", errors.Errorf("Can't find client actor")
	}

	author, err := ap.LoadActor(s.ctl.Storage, authIRI)
	if err != nil {
		s.logger.Errorf("Can't load client actor from db", err)
		return "", err
	}

	tags := make(vocab.ItemCollection, 0)

	objectsCollection := filters.ObjectsType.IRI(vocab.IRI(s.baseURL))
	allObjects, _ := s.ctl.Storage.Load(objectsCollection)
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
		Type: typ,
		// TODO(marius): when adding authentication to the command, we can set here the actor that executes it
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
		s.logger.Errorf("Can't save new actor", err)
		return "", err
	}

	//fmt.Printf("Added %q [%s]: %s\n", typ, name, newPerson.GetLink())
	/*
		if metaSaver, ok := s.ctl.Storage.(s.MetadataTyper); ok {
			if err := AddKeyToItem(metaSaver, newPerson, keyType); err != nil {
				Errf("Error saving metadata for %s: %s", name, err)
			}
		}
	*/
	return newPerson.GetLink(), nil
}
