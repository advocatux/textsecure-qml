package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *sqlx.DB
	dbDir  string
	dbFile string

	sessionsSchema = "CREATE TABLE IF NOT EXISTS sessions (id INTEGER PRIMARY KEY, name text, tel text, isgroup boolean, last string, timestamp integer, ctype integer)"
	sessionsInsert = "INSERT OR REPLACE INTO sessions (name, tel, isgroup, last, ctype) VALUES (:name, :tel, :isgroup, :last, :ctype)"
	sessionsSelect = "SELECT * FROM sessions"

	messagesSchema      = "CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY, sid integer, source text, message text, outgoing boolean, sentat integer, receivedat integer, ctype integer, attachment string, issent boolean, isread boolean)"
	messagesInsert      = "INSERT INTO messages (sid, source, message, outgoing, sentat, receivedat, ctype, attachment, issent, isread) VALUES (:sid, :source, :message, :outgoing, :sentat, :receivedat, :ctype, :attachment, :issent, :isread)"
	messagesSelectWhere = "SELECT * FROM messages WHERE sid = ?"

/*
	groupsSchema = "CREATE TABLE IF NOT EXISTS groups (id INTEGER PRIMARY KEY, groupid INTEGER, title TEXT, members TEXT, avatar BLOB, avatarid INTEGER, avatar_key BLOB, avatar_content_type TEXT, relay TEXT, active INTEGER DEFAULT 1)"
	groupsInsert = "INSERT INTO groups (groupid, title, members, active) VALUES (:groupid, :title, :members, :active)"
	groupsUpdate = "UPDATE groups SET members = :members, title = :title, WHERE groupid = :groupid"
*/
)

func setupDB() error {
	var err error

	dbDir = filepath.Join(dataDir, "db")
	dbFile = filepath.Join(dbDir, "db.sql")
	err = os.MkdirAll(dbDir, 0700)
	if err != nil {
		return err
	}

	db, err = sqlx.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	_, err = db.Exec(messagesSchema)
	if err != nil {
		return err
	}

	_, err = db.Exec(sessionsSchema)
	if err != nil {
		return err
	}
	/*
		_, err = db.Exec(groupsSchema)
		if err != nil {
			return err
		}
	*/
	return loadMessagesFromDB()
}

func saveSession(s *Session) error {
	res, err := db.NamedExec(sessionsInsert, s)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	s.ID = id
	return err
}

func saveMessage(m *Message) error {
	res, err := db.NamedExec(messagesInsert, m)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

func updateMessageSent(m *Message) error {
	_, err := db.NamedExec("UPDATE messages SET issent = :issent, sentat = :sentat WHERE id = :id", m)
	if err != nil {
		return err
	}
	return err
}

func updateMessageRead(m *Message) error {
	_, err := db.NamedExec("UPDATE messages SET isread = :isread WHERE id = :id", m)
	if err != nil {
		return err
	}
	return err
}

func updateSession(s *Session) error {
	_, err := db.NamedExec("UPDATE sessions SET timestamp = :timestamp, ctype = :ctype, last = :last WHERE id = :id", s)
	if err != nil {
		return err
	}
	return err
}

var allSessions []*Session

func loadMessagesFromDB() error {
	err := db.Select(&allSessions, sessionsSelect)
	if err != nil {
		return err
	}
	for _, s := range allSessions {
		s.When = humanize.Time(time.Unix(0, int64(1000000*s.Timestamp)))
		sessionsModel.sessions = append(sessionsModel.sessions, s)
		sessionsModel.Len++
		err = db.Select(&s.messages, messagesSelectWhere, s.ID)
		s.Len = len(s.messages)
		if err != nil {
			return err
		}
	}
	return nil
}
