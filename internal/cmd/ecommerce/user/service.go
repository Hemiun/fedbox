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

type UserService struct {
	db      common.Storage
	ctl     common.Control
	baseURL string
	logger  lw.Logger
}

func NewUserService(ctl common.Control, db common.Storage, baseURL string, l lw.Logger) *UserService {
	var target UserService
	target.db = db
	target.baseURL = baseURL
	target.logger = l
	return &target
}

const (
	//keyType    = fedbox.KeyTypeED25519
	clientName = "69c0d2c8-c105-45b2-bfbf-33ef8c7b770c"
)

func (s *UserService) clientUri() vocab.IRI {
	// get URI for main client
	//TODO:
	client, err := s.db.GetClient(clientName)
	if err != nil {
		s.logger.Errorf("Client not found. Unexpected error", err)
		return ""
	}

	if client == nil {
		s.logger.Errorf("Client not found", err)
		return ""
	}

	return vocab.IRI(client.GetId())
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

	author, err := ap.LoadActor(s.db, authIRI)
	if err != nil {
		s.logger.Errorf("Can't load client actor from db", err)
		return "", err
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
	//
	//	if metaSaver, ok := s.ctl.Storage.(s.MetadataTyper); ok {
	//		if err := AddKeyToItem(metaSaver, newPerson, keyType); err != nil {
	//			Errf("Error saving metadata for %s: %s", name, err)
	//		}
	//	}
	//
	return newPerson.GetLink(), nil
}
