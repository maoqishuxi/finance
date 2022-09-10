package writeToDatabase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTable(db *sql.DB, tableName string) {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS %s (
	    id INTEGER,
	    time TEXT PRIMARY KEY,
	    valueV REAL,
	    valueC REAL,
	    gain REAL
	)
`
	_, err := db.Exec(fmt.Sprintf(sqlStmt, tableName))
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("创建%s表成功\n", tableName)
}

func InsertData(db *sql.DB, ctx context.Context, data map[string][]byte) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println(err)
	}
	defer tx.Rollback()

	log.Println("开启事务")

	result, err := tx.ExecContext(ctx, "INSERT INTO file VALUES(?,?,?,?,?)",
		data["id"], data["time"], data["valueV"], data["valueC"], data["gain"])
	if err != nil {
		log.Println(err)
	}

	itemID, err := result.LastInsertId()
	fmt.Println(itemID)
	if err != nil {
		log.Println(err)
	}

	log.Println("插入数据完成，准备提交")

	if err = tx.Commit(); err != nil {
		log.Println(err)
	}

	log.Println("提交完成")
}
