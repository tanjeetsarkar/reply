package main

import (
	"database/sql"
	"log"
	"os"
)

const CreateConnectedUsersTable = `
CREATE TABLE IF NOT EXISTS connected_users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	status TEXT NOT NULL,
	connected TEXT NOT NULL
);`

type ConnectedUser struct {
	Name      string `json:"name"`
	Connected string `json:"connected"`
	Status    string `json:"status"`
}

type ConnectedUserDb struct {
	db *sql.DB
}

func NewConnectedUserRepo(db *sql.DB) *ConnectedUserDb {
	return &ConnectedUserDb{db: db}
}

func initializeDB(dbPath string, createScript string) (*sql.DB, error) {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(createScript)
		if err != nil {
			return nil, err
		}
		return db, nil
	} else if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (c *ConnectedUserDb) insert(connectedUser ConnectedUser) bool {
	_, err := c.db.Exec("INSERT INTO connected_users (name, status, connected) VALUES (?, ?, ?)", connectedUser.Name, connectedUser.Status, connectedUser.Connected)
	log.Println(err)
	return err == nil
}

func (c *ConnectedUserDb) update(connectedUser ConnectedUser) bool {
	_, err := c.db.Exec("UPDATE connected_users SET status = ?, connected = ? WHERE name = ?", connectedUser.Status, connectedUser.Connected, connectedUser.Name)
	log.Println(err)
	return err == nil
}

func (c *ConnectedUserDb) delete(connectedUser ConnectedUser) bool {
	_, err := c.db.Exec("DELETE FROM connected_users WHERE name = ?", connectedUser.Name)
	log.Println(err)
	return err == nil
}

func (c *ConnectedUserDb) get(name string) (ConnectedUser, bool) {
	var connectedUser ConnectedUser
	err := c.db.QueryRow("SELECT name, status, connected FROM connected_users WHERE name = ?", name).Scan(&connectedUser.Name, &connectedUser.Status, &connectedUser.Connected)
	if err != nil {
		return ConnectedUser{}, false
	}
	return connectedUser, true
}

func (c *ConnectedUserDb) getAll() []ConnectedUser {
	rows, err := c.db.Query("SELECT name, status, connected FROM connected_users")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var connectedUsers []ConnectedUser
	for rows.Next() {
		var connectedUser ConnectedUser
		err := rows.Scan(&connectedUser.Name, &connectedUser.Status, &connectedUser.Connected)
		if err != nil {
			return nil
		}
		connectedUsers = append(connectedUsers, connectedUser)
	}
	return connectedUsers
}

func (c *ConnectedUserDb) close() {
	c.db.Close()
}

func (c *ConnectedUserDb) deleteAll() bool {
	_, err := c.db.Exec("DELETE FROM connected_users")
	log.Println(err)
	return err == nil
}

func (c *ConnectedUserDb) count() int {
	var count int
	err := c.db.QueryRow("SELECT COUNT(*) FROM connected_users").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (c *ConnectedUserDb) exists(name string) bool {
	var count int
	err := c.db.QueryRow("SELECT COUNT(*) FROM connected_users WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (c *ConnectedUserDb) create() bool {
	_, err := c.db.Exec(CreateConnectedUsersTable)
	log.Println(err)
	return err == nil
}

func (c *ConnectedUserDb) drop() bool {
	_, err := c.db.Exec("DROP TABLE connected_users")
	log.Println(err)
	return err == nil
}

func (c *ConnectedUserDb) reset() bool {
	c.drop()
	return c.create()
}
