package store

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

const (
	KindCalendar = "Calendar"
	KindEntry    = "Entry"
)

// カレンダーのエントリー情報を格納する構造体。
type Entry struct {
	Year      int       `datastore:"-" json:"year"`
	Day       int       `datastore:"-" json:"day"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	Author    string    `json:"author"`
	Section   string    `json:"section"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// カレンダー Kind のキーを生成します。
func (s *Store) calendarKey(year int) *datastore.Key {
	return datastore.IDKey(KindCalendar, int64(year), nil)
}

// エントリー Kind のキーを生成します。
func (s *Store) entryKey(year, day int) *datastore.Key {
	calKey := s.calendarKey(year)
	return datastore.IDKey(KindEntry, int64(day), calKey)
}

// エントリーを取得します。
// エントリーが見つからない場合は nil を返します。
func (s *Store) GetEntry(ctx context.Context, year, day int) (*Entry, error) {
	entKey := s.entryKey(year, day)

	entry := new(Entry)
	err := s.client.Get(ctx, entKey, entry)
	if err == datastore.ErrNoSuchEntity {
		// データなし
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	entry.Year = year
	entry.Day = day

	return entry, nil
}

// year で指定した年のエントリーを取得します。
func (s *Store) GetEntries(ctx context.Context, year int) ([]*Entry, error) {
	calKey := s.calendarKey(year)

	entries := make([]*Entry, 0, 25)
	query := datastore.NewQuery(KindEntry).Ancestor(calKey)
	keys, err := s.client.GetAll(ctx, query, &entries)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		entries[i].Year = year
		entries[i].Day = int(key.ID)
	}

	return entries, nil
}

// エントリーを登録します。
func (s *Store) PutEntry(ctx context.Context, entry *Entry) error {
	now := time.Now()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	entry.UpdatedAt = now

	entKey := s.entryKey(entry.Year, entry.Day)

	_, err := s.client.Put(ctx, entKey, entry)

	return err
}
