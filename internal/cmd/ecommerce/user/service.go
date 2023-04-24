package user

import (
	"git.sr.ht/~mariusor/lw"
	_ "git.sr.ht/~mariusor/lw"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	ap "github.com/go-ap/fedbox/activitypub"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/metadata"
	"github.com/go-ap/filters"
	"time"
)

const (
	keyType = metadata.KeyTypeED25519
)

// UserService
type UserService struct {
	db        common.Storage
	ctl       common.Control
	baseURL   string
	logger    lw.Logger
	appActors vocab.Item
}

// NewUserService returns pointer to new UserService
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

// CreateUser creates new actor belonging to caller
// It's used for POST /user
func (s *UserService) CreateUser(caller vocab.Actor, ur UserDTO) (vocab.Item, error) {
	var it vocab.Item

	if ur.Name == "" || ur.Password == "" {
		return it, errors.Errorf("User credentials doesn't pass")
	}
	if !s.isSuperUser(caller) {
		err := errors.Errorf("Actor has insufficient privileges")
		return it, errors.NewForbidden(err, "Access denied")
	}

	author, err := ap.LoadActor(s.db, caller.GetLink())
	if err != nil {
		s.logger.Errorf("Can't load author from db", err)
		return it, err
	}

	tags := s.prepareTags(caller, ur.Tags)
	typ := vocab.PersonType

	now := time.Now().UTC()
	newPerson := &vocab.Person{
		Type:         typ,
		AttributedTo: author.GetLink(),
		Generator:    author.GetLink(),
		Published:    now,
		Summary: vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(ur.Comments)},
		},
		Updated: now,
		Name: vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(ur.Name)},
		},
		PreferredUsername: vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(ur.Name)},
		},
	}
	if len(tags) > 0 {
		newPerson.Tag = tags
	}

	if newPerson, err = s.ctl.AddActor(newPerson, []byte(ur.Password), &author); err != nil {
		s.logger.Errorf("Can't save new person", err)
		return it, err
	}

	if metaSaver, ok := s.db.(common.MetadataTyper); ok {
		err = func() error {
			if err := vocab.OnActor(newPerson, metadata.AddKeyToPerson(metaSaver, keyType)); err != nil {
				s.logger.Errorf("failed to process actor: %v", err)
				return err
			}
			if _, err := s.db.Save(newPerson); err != nil {
				s.logger.Errorf("can't save actor: %v", err)
				return err
			}
			return nil
		}()
		if err != nil {
			s.logger.Errorf("Error saving metadata for the actor: %v", err)
			return it, err
		}
	}
	return newPerson, nil
}

// DeleteUser delete actor that corresponds userID
// It's used for DELETE /user/{userID}
func (s *UserService) DeleteUser(caller vocab.Actor, userID string) error {
	// check access for the caller actor
	// You can use update endpoint if you are superuser or if you want upda self actor
	if !s.isSuperUser(caller) ||
		filters.ActorsType.IRI(vocab.IRI(s.baseURL)).AddPath(userID) == caller.ID {
		err := errors.Errorf("Actor has insufficient privileges")
		return errors.NewForbidden(err, "Access denied")
	}

	// Prepare actorIRI
	f := filters.FiltersNew()
	f.IRI = filters.ActorsType.IRI(vocab.IRI(s.baseURL)).AddPath(userID)

	// Load actor from repo
	actor, err := ap.LoadActor(s.db, f.GetLink())

	// It's important to check that actor is not empty. If your pass empty entity into db.Delete it will delete entire storage!
	if err != nil || !vocab.IsNotEmpty(actor) {
		s.logger.Errorf("Can't load actor from db", err)
		return errors.NewNotFound(err, "actor not found")
	}

	err = s.db.Delete(actor)
	if err != nil {
		s.logger.Errorf("Can't delete actor from db", err)
		return err
	}
	return nil
}

// FindUser returns UserDTO struct for actor that  corresponds userID
// It's used for GET /user/{userID}
func (s *UserService) FindUser(_ vocab.Actor, userID string) (*UserDTO, error) {
	// Available for all users without restrictions

	// Prepare actorIRI
	f := filters.FiltersNew()
	f.IRI = filters.ActorsType.IRI(vocab.IRI(s.baseURL)).AddPath(userID)

	// Load actor from repo
	actor, err := ap.LoadActor(s.db, f.GetLink())
	// It's important to check that actor is not empty. If your pass empty entity into db.Delete it will delete entire storage!
	if err != nil || !vocab.IsNotEmpty(actor) {
		s.logger.Errorf("Can't load actor from db", err)
		return nil, errors.NewNotFound(err, "actor not found")
	}

	tags := s.readItemTags(actor)

	res := UserDTO{
		Name:     actor.PreferredUsername.String(),
		Tags:     tags,
		Comments: actor.Summary.String(),
	}
	if err != nil {
		s.logger.Errorf("Can't delete actor from db", err)
		return nil, err
	}
	return &res, nil
}

// UpdateUser update actor that  corresponds userID
// It's used for PUT /user/{userID}
func (s *UserService) UpdateUser(caller vocab.Actor, userID string, ur UserDTO) (vocab.Item, error) {
	var it vocab.Item

	// You can use update endpoint if you are superuser or if you want update self actor
	if !s.isSuperUser(caller) ||
		filters.ActorsType.IRI(vocab.IRI(s.baseURL)).AddPath(userID) == caller.ID {
		err := errors.Errorf("Actor has insufficient privileges")
		return it, errors.NewForbidden(err, "Access denied")
	}

	author, err := ap.LoadActor(s.db, caller.GetLink())
	if err != nil {
		s.logger.Errorf("Can't load author from db", err)
		return it, err
	}

	f := filters.FiltersNew()
	f.IRI = filters.ActorsType.IRI(vocab.IRI(s.baseURL)).AddPath(userID)
	// Load actor from repo
	actor, err := ap.LoadActor(s.db, f.GetLink())
	// It's important to check that actor is not empty. If your pass empty entity into db.Delete it will delete entire storage!
	if err != nil || !vocab.IsNotEmpty(actor) {
		s.logger.Errorf("Can't load actor from db", err)
		return nil, errors.NewNotFound(err, "actor not found")
	}

	if ur.Name != "" {
		actor.Name = vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(ur.Name)},
		}
		actor.PreferredUsername = vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(ur.Name)},
		}
	}

	if ur.Comments != "" {
		actor.Summary = vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(ur.Comments)},
		}
	}

	tags := s.prepareTags(caller, ur.Tags)

	actor.AttributedTo = author.GetLink()
	actor.Updated = time.Now().UTC()

	if len(tags) > 0 {
		actor.Tag = tags
	}

	newItem, err := s.db.Save(actor)
	if err != nil {
		s.logger.Errorf("Can't update actor", err)
		return nil, errors.NewNotFound(err, "Can't update actor")
	}

	return newItem, nil
}

func (s *UserService) isSuperUser(actor vocab.Actor) bool {
	var flIsSuperUser bool
	err := vocab.OnCollectionIntf(s.appActors, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			err2 := vocab.OnObject(it, func(act *vocab.Object) error {
				if act.ID == actor.AttributedTo {
					flIsSuperUser = true
				}
				return nil
			})
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Errorf("can't check actor: %v", err)
		return false
	}
	return flIsSuperUser
}

func (s *UserService) readItemTags(obj vocab.Actor) []string {
	var res []string

	err := vocab.OnCollectionIntf(obj.Tag, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			res2 := vocab.OnObject(it, func(object *vocab.Object) error {
				res = append(res, object.Name.String())
				return nil
			})
			if res2 != nil {
				return res2
			}
		}
		return nil
	})
	if err != nil {
		s.logger.Errorf("can't read tags: %v", err)
		return nil
	}
	return res
}

func (s *UserService) prepareTags(owner vocab.Actor, src []string) vocab.ItemCollection {
	tags := make(vocab.ItemCollection, 0)

	existsTagFilter := filters.Filters{
		BaseURL:       vocab.IRI(s.baseURL),
		Authenticated: &owner,
		Type: filters.CompStrs{
			filters.CompStr{
				Str: string(vocab.ObjectType),
			},
		},
		IRI: vocab.IRI(s.baseURL + "/objects"),
		AttrTo: filters.CompStrs{
			filters.CompStr{
				Str: owner.ID.String(),
			},
		},
	}

	allObjects, _ := s.db.Load(existsTagFilter.GetLink())
	existsTags := map[string]*vocab.Object{}

	vocab.OnCollectionIntf(allObjects, func(col vocab.CollectionInterface) error {
		for _, it := range col.Collection() {
			vocab.OnObject(it, func(object *vocab.Object) error {
				existsTags[object.Name.First().Value.String()] = object
				return nil
			})
		}
		return nil
	})
	for _, t := range src {
		if el, ok := existsTags[t]; !ok {
			tag := vocab.ObjectNew("")
			tag.Name = vocab.NaturalLanguageValues{
				{Ref: vocab.NilLangRef, Value: vocab.Content(t)},
			}
			tag.AttributedTo = owner
			tags.Append(tag)
		} else {
			tags.Append(el)
		}
	}
	return tags
}
