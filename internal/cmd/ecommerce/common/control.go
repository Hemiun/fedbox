package common

import vocab "github.com/go-ap/activitypub"

type Control interface {
	AddActor(p *vocab.Person, pw []byte, author *vocab.Actor) (*vocab.Person, error)
}
