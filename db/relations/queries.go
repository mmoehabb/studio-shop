package relations

import (
	"fmt"

	anc "github.com/mmoehabb/studio-shop/ancillaries"
	"github.com/mmoehabb/studio-shop/db"
)

// retrieves an array of children of a specific parent
func GetSectionsOf(id int) ([]int, error) {
	conn := anc.Must(db.GetConnection()).(*db.Connection)
	rows := anc.Must(conn.Query("SELECT * FROM relations WHERE parent=$1", id)).([]any)
	var ids []int
	for _, row := range rows {
		ids = append(ids, parseRow(row.([]any)).Child)
	}
	return ids, nil
}

// inserts a new (parent-child) section relation in the database.
func Add(list []DataModel) error {
	conn := anc.Must(db.GetConnection()).(*db.Connection)
	query := "INSERT INTO relations VALUES "
	for _, data := range list {
		query += fmt.Sprintf("(%d,%d),", data.Parent, data.Child)
	}
	query = query[0 : len(query)-1]
	anc.Must(conn.Query(query))
	return nil
}

// removes a specific relation from the database
func Delete(data DataModel) error {
	conn := anc.Must(db.GetConnection()).(*db.Connection)
	anc.Must(conn.Query(
		"DELETE FROM relations WHERE parent=$1 AND child=$2",
		data.Parent,
		data.Child,
	))
	return nil
}

// removes specific parent relations from the database
func DeleteAll(parents []int) error {
	queryList := ""
	for _, id := range parents {
		queryList += fmt.Sprintf("%d,", id)
	}
	queryList = queryList[0 : len(queryList)-1]
	query := fmt.Sprintf("DELETE FROM relations WHERE parent IN (%s) OR child IN (%s)", queryList, queryList)

	conn := anc.Must(db.GetConnection()).(*db.Connection)
	anc.Must(conn.Query(query))
	return nil
}

// return true if there is no parent with the passed id
func IsAlbum(id int) bool {
	conn := anc.Must(db.GetConnection()).(*db.Connection)
	rows := anc.Must(conn.Query("SELECT child FROM relations WHERE parent=$1", id)).([]any)
	return len(rows) == 0
}
