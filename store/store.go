package store

import (
	"context"

	"cloud.google.com/go/datastore"
)

// データストアにアクセスするための情報を格納する構造体。
type Store struct {
	client *datastore.Client
}

// 新しい Store 構造体を生成します。
// ctx は App Engine のコンテキスト。
func NewStore(ctx context.Context, projectID string) (*Store, error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &Store{
		client: client,
	}, nil
}

func (s *Store) Close() error {
	return s.client.Close()
}
