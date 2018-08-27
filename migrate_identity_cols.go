package main

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/fatih/color"
	_ "gopkg.in/rana/ora.v4"
)

func main() {

	dataSourceNameArg := flag.String("dataSourceName", "ORCL", "Data source name from tnsnames.ora")
	userNameArg := flag.String("userName", "SYSTEM", "User name")
	passwordArg := flag.String("password", "ORACLE", "Password")

	flag.Parse()

	green := color.New(color.FgHiGreen).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	fmt.Fprintf(color.Output, "dataSourceName: %s\n", green(*dataSourceNameArg))
	fmt.Fprintf(color.Output, "userName: %s\n", green(*userNameArg))
	fmt.Fprintf(color.Output, "password: %s\n", green(*passwordArg))

	type sequence struct {
		owner       string
		tableName   string
		columnName  string
		dataDefault string
		nextVal     int64
	}
	var sequences []sequence
	var sequencesFinal []sequence

	db, err := sql.Open("ora", fmt.Sprintf("%v/%v@%v", *userNameArg, *passwordArg, *dataSourceNameArg))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Fprintf(color.Output, "Error connecting to the database: %s\n", red(err))
		return
	}

	rows, err := db.Query("SELECT table_name, column_name, data_default "+
		"FROM all_tab_columns "+
		"WHERE IDENTITY_column = :1 "+
		"AND data_default IS NOT null and owner = :2", "YES", *userNameArg)
	if err != nil {
		fmt.Fprintf(color.Output, "Error fetching data: %s\n", red(err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		var columnName string
		var dataDefault []byte
		rows.Scan(&tableName, &columnName, &dataDefault)
		s := sequence{owner: *userNameArg, tableName: tableName, columnName: columnName, dataDefault: string(dataDefault)}
		sequences = append(sequences, s)
	}

	for _, seq := range sequences {
		sql := fmt.Sprintf("SELECT %v as val from dual", seq.dataDefault)
		rows, err := db.Query(sql)
		if err != nil {
			fmt.Fprintf(color.Output, "Error fetching data: %s\n", red(err))
			return
		}
		defer rows.Close()

		var val int64
		for rows.Next() {
			rows.Scan(&val)
		}

		s := sequence{owner: seq.owner, tableName: seq.tableName, columnName: seq.columnName, dataDefault: seq.dataDefault, nextVal: val}
		sequencesFinal = append(sequencesFinal, s)
	}

	if len(sequencesFinal) > 0 {
		fmt.Println("SQL commands:")
		for _, seq := range sequencesFinal {
			alter := fmt.Sprintf("  ALTER TABLE %v.%v MODIFY %v GENERATED BY DEFAULT ON NULL AS IDENTITY (START WITH %d);", seq.owner, seq.tableName, seq.columnName, seq.nextVal)
			fmt.Fprintf(color.Output, "%s\n", magenta(alter))
		}
	}

}
