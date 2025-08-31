package pipeline

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var bucket = []byte("pipelines")

type Step struct {
	Type       string   `json:"type"`                  // es. "transcoder", "filewriter"
	FFmpegArgs []string `json:"ffmpeg_args,omitempty"` // args da passare a ffmpeg
	Path       string   `json:"path,omitempty"`        // solo per filewriter
}

type Pipeline struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

type Store struct {
	db *bolt.DB
}

func NewStore(path string) (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	// init bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Save(p Pipeline) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put([]byte(p.ID), data)
	})
}

func (s *Store) Get(id string) (*Pipeline, error) {
	var p Pipeline

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get([]byte(id))
		if v == nil {
			return fmt.Errorf("pipeline %s not found", id)
		}
		return json.Unmarshal(v, &p)
	})
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *Store) List() ([]Pipeline, error) {
	var result []Pipeline

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.ForEach(func(k, v []byte) error {
			var p Pipeline
			if err := json.Unmarshal(v, &p); err != nil {
				return err
			}
			result = append(result, p)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
