package repository

import "database/sql"

// The connection manager to handle the read-write connection

func GetConnection(readonly bool) (*sql.DB, error) {
	db, err := sql.Open("mysql", "dbrpmt:DbRESTFu1-pmt@(db.91qpzs.net:3326)/pmt?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		return nil, err
	}

	return db, nil
}
