package store

import "google.golang.org/appengine/datastore"

const (
	KindCalendar = "Calendar"
	KindEntry    = "Entry"
)

// カレンダーのエントリー情報を格納する構造体。
type Entry struct {
	Year    int `datastore:"-"`
	Day     int `datastore:"-"`
	Title   string
	Url     string
	Author  string
	Section string
}

// カレンダー Kind のキーを生成します。
func (s *Store) calendarKey(year int) *datastore.Key {
	return datastore.NewKey(s.ctx, KindCalendar, "", int64(year), nil)
}

// エントリー Kind のキーを生成します。
func (s *Store) entryKey(year, day int) *datastore.Key {
	calKey := s.calendarKey(year)
	return datastore.NewKey(s.ctx, KindEntry, "", int64(day), calKey)
}

// エントリーを取得します。
// エントリーが見つからない場合は nil を返します。
func (s *Store) GetEntry(year, day int) (*Entry, error) {
	entKey := s.entryKey(year, day)

	entry := new(Entry)
	err := datastore.Get(s.ctx, entKey, entry)
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
func (s *Store) GetEntries(year int) ([]*Entry, error) {
	calKey := s.calendarKey(year)

	entries := make([]*Entry, 0, 25)
	query := datastore.NewQuery(KindEntry).Ancestor(calKey)
	keys, err := query.GetAll(s.ctx, &entries)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		entries[i].Year = year
		entries[i].Day = int(key.IntID())
	}

	return entries, nil
}

// エントリーを登録します。
func (s *Store) PutEntry(entry *Entry) error {
	entKey := s.entryKey(entry.Year, entry.Day)

	_, err := datastore.Put(s.ctx, entKey, entry)

	return err
}
