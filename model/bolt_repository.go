package model

import (
	bolt "github.com/coreos/bbolt"
	"errors"
	"database/sql"
)

const (
	PROVIDER_BUCKET string = "Providers"
)

type BBoltRepository struct {
	dbFile string
}

func CreateRepository(fileName string) (BBoltRepository, error) {
	repo := BBoltRepository{fileName}

	err := repo.ExecuteTx(nil, func(db *bolt.DB) error {
		return db.Update(func (tx *bolt.Tx) error {
			_, err2 := tx.CreateBucketIfNotExists([]byte(PROVIDER_BUCKET))
			return err2
		})
	})

	return repo, err
}

func (r BBoltRepository) ExecuteTx(options *bolt.Options, f func(db *bolt.DB) error) error {
	db, err := bolt.Open(r.dbFile, 0666, options)
	if err != nil {
		return err
	}
	defer db.Close()
	return f(db)
}

func (r BBoltRepository) ExecuteOperation(ft func (fn func(tx *bolt.Tx) error) error, fb func(b *bolt.Bucket) error ) error {
	return ft(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PROVIDER_BUCKET))
		if b != nil {
			return fb(b)
		} else {
			return errors.New("Error getting providers bucket")
		}
	})
}

func (r BBoltRepository) ExecuteReadOperation(f func(b *bolt.Bucket) error) error {
	return r.ExecuteTx(&bolt.Options{ReadOnly: true}, func (db *bolt.DB) error {
		return r.ExecuteOperation(db.View, f)
	})
}

func (r BBoltRepository) ExecuteWriteOperation(f func(b *bolt.Bucket) error) error {
	return r.ExecuteTx(nil, func (db *bolt.DB) error {
		return r.ExecuteOperation(db.Update, f)
	})
}

func (r BBoltRepository) GetAllProviders() (Providers, error) {

	var providers Providers

	err := r.ExecuteReadOperation(func (b *bolt.Bucket) error {
		b.ForEach(func(k, v []byte) error {
			providers = append(providers, Provider{string(k), string(v)})
			return nil
		})

		return nil
	})
	return providers, err
}

func (r BBoltRepository) GetProvider(id string) (*Provider, error) {
	var provider *Provider = nil
	err := r.ExecuteReadOperation(func(b *bolt.Bucket) error {
		value := b.Get([]byte(id))
		if value != nil {
			provider = &Provider{id, string(value)}
		} else {
			return sql.ErrNoRows
		}
		return nil
	})
	return provider, err
}

func (r BBoltRepository) CreateProvider(provider *Provider) (*Provider, error) {
	//log.Info("Trying to create provider " + provider.Id)
	return provider, r.ExecuteWriteOperation(func(b *bolt.Bucket) error {
		existing := b.Get([]byte(provider.Id))
		if existing != nil {
			return sql.ErrNoRows
		} else {
			return b.Put([]byte(provider.Id), []byte(provider.Name))
		}
	})
}
