package main

import (
	"fmt"

	"gopkg.in/rana/ora.v4"
)

func main() {

	// example usage of the ora package driver
	// connect to a server and open a session
	env, err := ora.OpenEnv()
	defer env.Close()
	if err != nil {
		panic(err)
	}
	srvCfg := ora.SrvCfg{Dblink: "ZISD"}
	srv, err := env.OpenSrv(srvCfg)
	defer srv.Close()
	if err != nil {
		panic(err)
	}
	sesCfg := ora.SesCfg{
		Username: "nxt",
		Password: "cetin1",
	}
	ses, err := srv.OpenSes(sesCfg)
	defer ses.Close()
	if err != nil {
		panic(err)
	}

	schemas := []string{"NETCAT", "PUBLIC_VIEW", "NXT"}
	for _, schema := range schemas {
		fmt.Println(">> schema:", schema)

		// Drop table
		_, err = ses.PrepAndExe("drop table tmp_ident_col")

		// Crate table
		_, err = ses.PrepAndExe(fmt.Sprintf("create table tmp_ident_col as SELECT table_name, column_name, to_lob(data_default) as data_default "+
			"FROM dba_tab_columns WHERE IDENTITY_column = '%v' AND data_default IS NOT null and owner = '%v'", "YES", schema))
		if err != nil {
			panic(err)
		}

		// Fetch records
		stmtQry, err := ses.Prep("SELECT table_name, column_name, data_default from tmp_ident_col")
		defer stmtQry.Close()
		if err != nil {
			panic(err)
		}
		rset, err := stmtQry.Qry()
		if err != nil {
			panic(err)
		}
		for rset.Next() {
			fmt.Println("tableName:", rset.Row[0], "columnName:", rset.Row[1], "dataDefault:", rset.Row[2])
		}
		if err := rset.Err(); err != nil {
			panic(err)
		}
	}
}
