package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/buntdb"

	"git.mkz.me/mycroft/asoai/internal/session"
)

type DB struct {
	handle *buntdb.DB
}

// Opens the database located in the given file path
func Open(filePath string) (*DB, error) {
	db, err := buntdb.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err)
	}

	err = db.CreateIndex("sessions", "session:*", buntdb.IndexString)
	if err != nil {
		return nil, fmt.Errorf("could not create index: %v", err)
	}

	return &DB{
		handle: db,
	}, nil
}

// Opens the database or fail and exits on error
func OpenOrFail(directory string) *DB {
	db, err := Open(directory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return db
}

// Save session in database
func (db *DB) SetSession(name string, session session.Session) error {
	encoded, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return db.handle.Update(func(tx *buntdb.Tx) error {
		_, _, err = tx.Set(fmt.Sprintf("session:%s", name), string(encoded), nil)
		return err
	})
}

// Retrieve session from database
func (db *DB) GetSession(name string) (session.Session, error) {
	var session session.Session
	var val string
	var err error

	err = db.handle.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(fmt.Sprintf("session:%s", name))
		return err
	})

	if err != nil {
		return session, fmt.Errorf("could not retrieve session: %v", err)
	}

	err = json.Unmarshal([]byte(val), &session)
	if err != nil {
		return session, fmt.Errorf("could not unmarshal session: %v", err)
	}

	return session, err
}

// List sessions from database and returns an array of strings
func (db *DB) ListSessions() ([]string, error) {
	var sessions []string

	err := db.handle.View(func(tx *buntdb.Tx) error {
		tx.Ascend("sessions", func(key, val string) bool {
			sessions = append(sessions, strings.Split(key, ":")[1])
			return true
		})
		return nil
	})

	return sessions, err

}

// Set current session in database
func (db *DB) SetCurrentSession(name string) error {
	err := db.handle.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("current", fmt.Sprintf("session:%s", name), nil)
		return err
	})

	return err
}

// Get current session from database
func (db *DB) GetCurrentSession() (string, error) {
	var name string

	err := db.handle.View(func(tx *buntdb.Tx) error {
		var err error
		name, err = tx.Get("current")
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		return nil
	})

	// Not having "session:" has a prefix should not be possible
	if !strings.HasPrefix(name, "session:") {
		return name, err
	}

	return strings.Split(name, ":")[1], err
}

// Delete given session in database
func (db *DB) DeleteSession(name string) error {
	err := db.handle.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(fmt.Sprintf("session:%s", name))
		return err
	})

	return err
}

// Shrink/compact database
func (db *DB) Shrink() error {
	err := db.handle.Shrink()
	if err != nil {
		return fmt.Errorf("could not shrink database: %v", err)
	}
	return err
}

// Close the database handle
func (db *DB) Close() error {
	err := db.handle.Close()
	if err != nil {
		return fmt.Errorf("could not close database: %v", err)
	}
	return err
}
