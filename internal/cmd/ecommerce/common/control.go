package common

import vocab "github.com/go-ap/activitypub"

// Control - is interface for access to  cmd.Control fedbox structure
type Control interface {
	AddActor(p *vocab.Person, pw []byte, author *vocab.Actor) (*vocab.Person, error)
}
