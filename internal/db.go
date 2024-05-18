package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/buntdb"
)

func OpenDB() (*buntdb.DB, error) {
	db, err := buntdb.Open("data.db")
	if err != nil {
		return nil, fmt.Errorf("could not open database: %v", err)
	}

	err = db.CreateIndex("sessions", "session:*", buntdb.IndexString)
	if err != nil {
		return nil, fmt.Errorf("could not create index: %v", err)
	}

	return db, nil
}

func DbSetSession(sessionUuid string, session Session) error {
	db, err := OpenDB()
	if err != nil {
		return fmt.Errorf("could not open db: %v", err)
	}
	defer db.Close()

	encoded, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return db.Update(func(tx *buntdb.Tx) error {
		tx.Set(fmt.Sprintf("session:%s", sessionUuid), string(encoded), nil)
		return nil
	})
}

func DbCreateSession(sessionUuid string, session Session) error {
	return DbSetSession(sessionUuid, session)
}

func DbGetSession(sessionUuid string) (Session, error) {
	var val string
	var session Session

	db, err := OpenDB()
	if err != nil {
		return session, fmt.Errorf("could not open db: %v", err)
	}
	defer db.Close()

	err = db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(fmt.Sprintf("session:%s", sessionUuid))
		return err
	})

	err = json.Unmarshal([]byte(val), &session)

	return session, err
}

func DbListSessions() ([]string, error) {
	var sessions []string

	db, err := OpenDB()
	if err != nil {
		return []string{}, fmt.Errorf("could not open database: %v", err)
	}
	defer db.Close()

	err = db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("sessions", func(key, val string) bool {
			sessions = append(sessions, strings.Split(key, ":")[1])
			return true
		})
		return nil
	})

	return sessions, err
}

func DbGetCurrentSession() (string, error) {
	var uuid string

	db, err := OpenDB()
	if err != nil {
		return uuid, fmt.Errorf("could not open database: %v", err)
	}
	defer db.Close()

	err = db.View(func(tx *buntdb.Tx) error {
		var err error
		uuid, err = tx.Get("current")
		if err != nil && err != buntdb.ErrNotFound {
			return err
		}
		return nil
	})

	return uuid, err
}

func DbSetCurrentSession(session string) error {
	db, err := OpenDB()
	if err != nil {
		return fmt.Errorf("could not open database: %v", err)
	}
	defer db.Close()

	err = db.Update(func(tx *buntdb.Tx) error {
		_, _, err = tx.Set("current", session, nil)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
