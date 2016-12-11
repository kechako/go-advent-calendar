package store

import "golang.org/x/net/context"

// データストアにアクセスするための情報を格納する構造体。
type Store struct {
	ctx context.Context
}

// 新しい Store 構造体を生成します。
// ctx は App Engine のコンテキスト。
func NewStore(ctx context.Context) *Store {
	return &Store{
		ctx: ctx,
	}
}
