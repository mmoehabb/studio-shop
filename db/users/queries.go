package users

import (
	"errors"

	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db"
)

func Add(username, password string) error {
	conn := anc.Must(db.GetConnection()).(*db.Connection)

	res := anc.Must(conn.SeqQuery("SELECT * FROM users WHERE username=$1", username)).([]any)
	if len(res) != 0 {
		conn.Close()
		return errors.New("username already found.")
	}

	anc.Must(conn.Query("INSERT INTO users VALUES ($1, $2)", username, password))
	return nil
}

func Get(username string) (DataModel, error) {
	conn := anc.Must(db.GetConnection()).(*db.Connection)

	res := anc.Must(conn.Query("SELECT * FROM users WHERE username=$1", username)).([]any)
	if len(res) == 0 {
		return DataModel{}, errors.New("couldn't find username.")
	}
	user := parseRow(res[0].([]any))
	return user, nil
}

// returns true if there no data in users table
func IsEmpty() bool {
	conn := anc.Must(db.GetConnection()).(*db.Connection)
	res := anc.Must(conn.Query("SELECT * FROM users")).([]any)
	return len(res) == 0
}
