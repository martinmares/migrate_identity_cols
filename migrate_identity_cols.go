package main

import (
	"database/sql"
	"fmt"

	_ "gopkg.in/rana/ora.v4"
)

func main() {
	db, err := sql.Open("ora", "nxt/cetin1@ZISD")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		return
	}

	schemas := []string{"KFA", "CIP", "NETCAT", "IFACE", "NXT", "PUBLIC_VIEW"}
	for _, schema := range schemas {

		fmt.Println("schema:", schema)

		rows, err := db.Query("SELECT table_name, column_name, data_default "+
			"FROM dba_tab_columns "+
			"WHERE IDENTITY_column = :1 "+
			"AND data_default IS NOT null and owner = :2", "YES", schema)

		if err != nil {
			fmt.Println("Error fetching data.")
			fmt.Println(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var tableName string
			var columnName string
			var dataDefault []byte
			rows.Scan(&tableName, &columnName, &dataDefault)
			fmt.Printf("Table name: %v, column: %v, data: %v\n", tableName, columnName, string(dataDefault))
		}
	}
}
