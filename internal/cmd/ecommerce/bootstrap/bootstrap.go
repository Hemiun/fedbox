package bootstrap

import (
	"fmt"
	vocab "github.com/go-ap/activitypub"
	"github.com/go-ap/errors"
	ap "github.com/go-ap/fedbox/activitypub"
	"github.com/go-ap/fedbox/internal/cmd/ecommerce/common"
	"github.com/urfave/cli/v2"
	"net/url"
	"path"
	"time"
)

const superuserName = "superuser"
const superuserDefaultSecret = "admin"

func InitNewApp(c *cli.Context, ctl common.Control, baseURL string, db common.Storage) error {
	//create oauth client
	clientSecret := []byte(c.String("secret"))

	redirectURIs := c.StringSlice("redirectUri")
	if len(redirectURIs) < 1 {
		return errors.Newf("Need to provide redirect URI for the client")
	}
	var appURL vocab.IRI
	for i, redirectUrl := range redirectURIs {
		if u, err := url.ParseRequestURI(redirectUrl); err == nil {
			u.Path = path.Clean(u.Path)
			redirectURIs[i] = u.String()
			if appURL == "" {
				appURL = vocab.IRI(u.String())
			}
		}
	}

	clientID, err := ctl.AddClient(clientSecret, redirectURIs, nil)
	if err != nil {
		return err
	}

	appActorUri := vocab.IRI(fmt.Sprintf("%s/actors/%s", baseURL, clientID))

	author, err := ap.LoadActor(db, appActorUri)
	if err != nil {
		return err
	}
	//appActor := vocab.PersonNew(appActorUri)

	superuserSecret := []byte(superuserDefaultSecret)
	now := time.Now().UTC()
	adminActor := &vocab.Person{
		Type:         vocab.PersonType,
		AttributedTo: author.GetLink(),
		Generator:    author.GetLink(),
		Published:    now,
		Updated:      now,
		Name: vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(superuserName)},
		},
		PreferredUsername: vocab.NaturalLanguageValues{
			{Ref: vocab.NilLangRef, Value: vocab.Content(superuserName)},
		},
		URL: appURL,
	}

	adminActor, err = ctl.AddActor(adminActor, superuserSecret, &author)
	if err != nil {
		return err
	}
	printReport(clientID, adminActor.ID.String(), string(superuserSecret))
	return nil
}

func printReport(clientID string, superuser string, superuserSecret string) {
	template := `
----------------------------------------------------------------------------------------------
Your application has successfully initialized
Be sure to save the following data:
oAuth ClientID: %s
superuser: %s
superuser password: %s
You must change the password for the user after the first login!
------------------------------------------------------------------------------------------------
`
	fmt.Printf(template, clientID, superuser, superuserSecret)

}
