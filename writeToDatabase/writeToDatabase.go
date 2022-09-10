package writeToDatabase

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	Id     int
	Time   []byte
	ValueV []byte
	ValueC []byte
	Gain   []byte
}

func CreateTable(db *sql.DB, tableName string) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS %s (
	    id INTEGER,
	    time TEXT PRIMARY KEY,
	    valueV TEXT,
	    valueC TEXT,
	    gain TEXT
	)
`
	_, err := db.Exec(fmt.Sprintf(sqlStmt, tableName))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("创建%s表成功\n", tableName)
}

func QueryData(db *sql.DB, num int) []Item {
	rows, err := db.Query("select * from file limit ?", num)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	result := make([]Item, 0)
	var ret Item
	for rows.Next() {
		if err := rows.Scan(&ret.Id, &ret.Time, &ret.ValueV, &ret.ValueC, &ret.Gain); err != nil {
			log.Printf("query error: %v\n", err)
		}
		result = append(result, ret)
	}

	return result
}

func QueryID(db *sql.DB) int {
	var id int
	row := db.QueryRow("SELECT id FROM file ORDER BY id desc LIMIT 1")
	if err := row.Scan(&id); err != nil {
		log.Printf("QueryID: %v\n", err)
	}

	return id
}

func InsertData(db *sql.DB, data Item) {
	result, err := db.Exec("INSERT INTO file VALUES (?,?,?,?,?)",
		data.Id, data.Time, data.ValueV, data.ValueC, data.Gain)
	if err != nil {
		log.Println(err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(id)
}
