package cryptoname

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
	"time"
)

type storage struct {
	db *badger.DB
}

func OpenStorage(path string) (*storage, error) {
	db, err := badger.Open(badger.DefaultOptions(path).
		WithInMemory(path == "").
		WithLogger(nil).
		WithCompression(options.None).
		WithSyncWrites(false).
		WithMetricsEnabled(false),
	)
	if err != nil {
		return nil, err
	}
	return &storage{db}, nil
}

func (st *storage) Close() error {
	return st.db.Close()
}
func (st *storage) Load(record *Record) error {
	tx := st.db.NewTransaction(false)
	defer tx.Discard()
	item, err := tx.Get(record.PK)
	if err != nil {
		return err
	}
	if item.IsDeletedOrExpired() {
		return badger.ErrKeyNotFound
	}
	val, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}
	return record.Decode(val)
}

var ErrOldRecord = errors.New("old record")
var ErrDuplicateRecord = errors.New("duplicate record")

func (st *storage) GetVersion(record *Record) (uint64, error) {
	tx := st.db.NewTransaction(false)
	defer tx.Discard()
	item, err := tx.Get(record.PK)
	if err != nil {
		return 0, err
	}
	if item.IsDeletedOrExpired() {
		return 0, badger.ErrKeyNotFound
	}
	return item.Version(), nil
}
func (st *storage) Store(record *Record, ttl time.Duration) error {
	tx := st.db.NewTransaction(true)
	defer tx.Discard()
	item, err := tx.Get(record.PK)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return err
		}
	} else {
		if !item.IsDeletedOrExpired() {
			val, err := item.ValueCopy(nil)
			if err == nil {
				found := new(Record)
				err := found.Decode(val)
				if err == nil {
					if found.Rev == record.Rev {
						return ErrDuplicateRecord
					}
					if found.Rev > record.Rev {
						return ErrOldRecord
					}
				}
			}
		}
	}
	err = tx.SetEntry(badger.NewEntry(record.PK, record.Encode()).WithTTL(ttl))
	if err != nil {
		return err
	}
	return tx.Commit()
}
