//go:build storage_all || (!storage_boltdb && !storage_fs && !storage_badger && !storage_sqlite)

package fedbox

import (
	"git.sr.ht/~mariusor/lw"
	authbadger "github.com/go-ap/auth/badger"
	authsqlite "github.com/go-ap/auth/sqlite"
	"github.com/go-ap/errors"
	"github.com/go-ap/fedbox/internal/config"
	"github.com/go-ap/processing"
	"github.com/go-ap/storage-badger"
	"github.com/go-ap/storage-boltdb"
	fs "github.com/go-ap/storage-fs"
	sqlite "github.com/go-ap/storage-sqlite"
	"github.com/openshift/osin"
)

func getBadgerStorage(c config.Options, l lw.Logger) (processing.Store, osin.Storage, error) {
	path := c.BaseStoragePath()
	conf := badger.Config{Path: path, Logger: l}
	if l != nil {
		l.Debugf("Initializing badger storage at %s", path)
	}
	db, err := badger.New(conf)
	if err != nil {
		return db, nil, err
	}
	authConf := authbadger.Config{Path: c.BadgerOAuth2(path)}
	oauth := authbadger.New(authConf)
	return db, oauth, nil
}

func getBoltStorage(c config.Options, l lw.Logger) (processing.Store, osin.Storage, error) {
	path := c.BaseStoragePath()
	l = l.WithContext(lw.Ctx{"path": path})
	l.Debugf("Initializing boltdb storage")
	db, err := boltdb.New(boltdb.Config{
		Path:    path,
		BaseURL: c.BaseURL,
		LogFn:   l.Debugf,
		ErrFn:   l.Warnf,
	})
	if err != nil {
		return nil, nil, err
	}
	return db, db, nil
}

func getFsStorage(c config.Options, l lw.Logger) (processing.Store, osin.Storage, error) {
	p := c.BaseStoragePath()
	l = l.WithContext(lw.Ctx{"path": p})
	l.Debugf("Initializing fs storage")
	db, err := fs.New(fs.Config{
		Path: p,
		CacheEnable: c.StorageCache,
		LogFn: l.Debugf,
		ErrFn: l.Warnf,
	})
	if err != nil {
		return nil, nil, err
	}
	return db, db, nil
}

func getSqliteStorage(c config.Options, l lw.Logger) (processing.Store, osin.Storage, error) {
	path := c.BaseStoragePath()
	l.Debugf("Initializing sqlite storage at %s", path)
	oauth := authsqlite.New(authsqlite.Config{
		Path:  path,
		LogFn: InfoLogFn(l),
		ErrFn: ErrLogFn(l),
	})
	db, err := sqlite.New(sqlite.Config{
		Path:        path,
		CacheEnable: c.StorageCache,
	})
	if err != nil {
		return nil, nil, errors.Annotatef(err, "unable to connect to sqlite storage")
	}
	return db, oauth, nil
}

func Storage(c config.Options, l lw.Logger) (processing.Store, osin.Storage, error) {
	switch c.Storage {
	case config.StorageBoltDB:
		return getBoltStorage(c, l)
	case config.StorageBadger:
		return getBadgerStorage(c, l)
	case config.StorageSqlite:
		return getSqliteStorage(c, l)
	case config.StorageFS:
		return getFsStorage(c, l)
	}
	return nil, nil, errors.NotImplementedf("Invalid storage type %s", c.Storage)
}
