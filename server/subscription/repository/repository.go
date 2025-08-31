package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/data"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/domain"
	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("subscriptions")

type Repository struct {
	db *bolt.DB
}

// Delete implements domain.Repository.
func (r *Repository) Delete(ctx context.Context, id string) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Delete([]byte(id))
	})
}

// GetCursor implements domain.Repository.
func (s *Repository) GetCursor(ctx context.Context, id string) (int64, error) {
	var cursor int64

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("subscriptions"))
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("subscription %s not found", id)
		}

		var data struct {
			Cursor int64 `json:"cursor"`
		}

		if err := json.Unmarshal(v, &data); err != nil {
			return err
		}
		cursor = data.Cursor
		return nil
	})

	if err != nil {
		return -1, err
	}

	return cursor, nil
}

// List implements domain.Repository.
func (r *Repository) List(ctx context.Context, start int64, limit int) (*[]data.Subscription, error) {
	var subs []data.Subscription

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.ForEach(func(k, v []byte) error {
			var sub data.Subscription
			if err := json.Unmarshal(v, &sub); err != nil {
				return err
			}
			subs = append(subs, sub)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return &subs, nil
}

// Submit implements domain.Repository.
func (s *Repository) Submit(ctx context.Context, sub *data.Subscription) (*data.Subscription, error) {
	if sub.Id == "" {
		sub.Id = uuid.NewString()
	}

	data, err := json.Marshal(sub)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("subscriptions"))
		return b.Put([]byte(sub.Id), data)
	})

	if err != nil {
		return nil, err
	}

	return sub, nil
}

// UpdateByExample implements domain.Repository.
func (s *Repository) UpdateByExample(ctx context.Context, example *data.Subscription) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("subscriptions"))

		return b.ForEach(func(k, v []byte) error {
			var sub data.Subscription
			if err := json.Unmarshal(v, &sub); err != nil {
				return err
			}

			if sub.Id == example.Id || sub.URL == example.URL {
				// aggiorna i campi
				sub.URL = example.URL
				sub.Params = example.Params
				sub.CronExpr = example.CronExpr

				data, err := json.Marshal(sub)
				if err != nil {
					return err
				}

				if err := b.Put(k, data); err != nil {
					return err
				}
			}

			return nil
		})
	})
}

func New(db *bolt.DB) domain.Repository {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})

	return &Repository{
		db: db,
	}
}
