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
		log.Println("createTable error: ", err)
		return
	}

	fmt.Printf("创建%s表成功\n", tableName)
}

func QueryData(db *sql.DB, tableName string, num int) []Item {
	rows, err := db.Query(fmt.Sprintf("select * from %s limit ?", tableName), num)
	if err != nil {
		log.Println("QueryData error: ", err)
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

func QueryID(db *sql.DB, tableName string) int64 {
	var id int64
	row := db.QueryRow(fmt.Sprintf("SELECT id FROM %s ORDER BY id desc LIMIT 1", tableName))
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0
		}
		log.Println("QueryID error", err)
	}

	return id
}

func InsertData(db *sql.DB, tableName string, data Item, num int64) int64 {
	result, err := db.Exec(fmt.Sprintf("INSERT INTO %s VALUES (?,?,?,?,?)", tableName),
		num, data.Time, data.ValueV, data.ValueC, data.Gain)
	if err != nil {
		log.Println("InsertData error: ", err)
		return -1
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
	}

	return id
}
