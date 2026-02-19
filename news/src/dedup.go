package src

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type DedupDB struct {
	db *sql.DB
}

func OpenDedupDB(path string) (*DedupDB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("could not open dedup database: %w", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posted_items (
		guid TEXT PRIMARY KEY,
		feed_url TEXT NOT NULL,
		title TEXT NOT NULL DEFAULT '',
		tx_hash TEXT NOT NULL DEFAULT '',
		posted_at INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not create posted_items table: %w", err)
	}
	return &DedupDB{db: db}, nil
}

func (d *DedupDB) CleanOld(days int) error {
	cutoff := time.Now().Unix() - int64(days*86400)
	_, err := d.db.Exec("DELETE FROM posted_items WHERE posted_at < ?", cutoff)
	return err
}
func (d *DedupDB) Close() error {
	return d.db.Close()
}
func (d *DedupDB) IsPosted(guid string) (bool, error) {
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM posted_items WHERE guid = ?", guid).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (d *DedupDB) MarkPosted(guid, feedUrl, title, txHash string) error {
	_, err := d.db.Exec(
		"INSERT OR IGNORE INTO posted_items (guid, feed_url, title, tx_hash, posted_at) VALUES (?, ?, ?, ?, ?)",
		guid, feedUrl, title, txHash, time.Now().Unix(),
	)
	return err
}
