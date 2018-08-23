package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-oci8"
)

func main() {
	db, err := sql.Open("oci8", "nxt/cetin1@ZISD")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Error connecting to the database: %s\n", err)
		return
	}

	// rows, err := db.Query("SELECT table_name, column_name FROM dba_tab_columns WHERE IDENTITY_column = :1 AND data_default IS NOT null and owner = :2", "YES", "NETCAT")
	// if err != nil {
	// 	fmt.Println("Error fetching addition")
	// 	fmt.Println(err)
	// 	return
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var tableName string
	// 	var columnName string
	// 	rows.Scan(&tableName, &columnName)
	// 	fmt.Printf("Table name: %v, column: %v\n", tableName, columnName)
	// }

	schemas := []string{"NETCAT", "PUBLIC_VIEW"}
	for _, schema := range schemas {

		fmt.Println("schema:", schema)

		_, err = db.Exec("drop table tmp_ident_col")
		_, err = db.Exec(fmt.Sprintf("create table tmp_ident_col as SELECT table_name, column_name, to_lob(data_default) as data_default FROM dba_tab_columns WHERE IDENTITY_column = '%v' AND data_default IS NOT null and owner = '%v'", "YES", schema))
		// ORA-01036: nepřípustné jméno nebo číslo proměnné
		//_, err = db.Exec("create table tmp_ident_col as SELECT table_name, column_name, to_lob(data_default) as data_default FROM dba_tab_columns WHERE identity_column = :1 AND data_default IS NOT null and owner = :2 ", "YES", schema)
		if err != nil {
			fmt.Println(err)
			return
		}

		//rows, err := db.Query("SELECT table_name, column_name, data_default FROM dba_tab_columns WHERE IDENTITY_column = :1 AND data_default IS NOT null and owner = :2", "YES", schema)
		rows, err := db.Query("SELECT table_name, column_name, data_default from tmp_ident_col")

		if err != nil {
			fmt.Println("Error fetching data")
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
		_, err = db.Exec("drop table tmp_ident_col")
	}
}
