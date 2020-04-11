package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

type Repository interface {
	Invoke(csiName string, data interface{}, db *sql.DB, transaction *sql.Tx) (interface{}, error)
}

// func QuerySql(sql string) ([]map[string]interface{}, error) {
// 	conn, _ := dbr.Open("mysql", "dbrpmt:DbRESTFu1-pmt@(db.91qpzs.net:3326)/pmt?charset=utf8&parseTime=True&loc=Local", nil)
// 	conn.SetMaxOpenonns(10)

// 	defer conn.Close()

// 	// create a sesson for each business unit of execution (e.g. a web request or goworkers job)
// sess := conn.NewSession(nil)

// 	result := []map[string]interface{}{}
// 	count, err := sess.electBySql(sql).Load(&result)

// 	fmt.Printf("get %v records.", count)

// return result, err
// }

func Invoke(query string, args ...interface{}) (interface{}, error) {
	db, err := GetConnection(false)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	var objects []map[string]interface{}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	columns, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		dataRow := map[string]interface{}{}

		for i, value := range values {
			read(value, dataRow, columns[i])
		}

		objects = append(objects, dataRow)
	}

	return objects, err
}

// read row to data
func read(data interface{}, row map[string]interface{}, column *sql.ColumnType) error {
	columnName := column.Name()
	dbType := column.DatabaseTypeName()

	switch dbType {
	case "INT", "TINYINT", "UINT", "UTINYINT", "MEDIUMINT":
		dest := new(int)
		ok := ConvertAssign(&dest, data)
		if ok == nil {
			row[columnName] = dest
		}
		return ok
	case "BIT":
		v, ok := data.([]byte)
		if ok && len(v) == 1 {
			row[columnName] = v[0] == 1
		}
		return nil
	case "NullInt64":
		dest := new(int64)
		ok := ConvertAssign(&dest, data)
		if ok == nil {
			row[columnName] = dest
		}
		return ok
	case "NullTime", "DATETIME":
		dest := new(time.Time)
		ok := ConvertAssign(&dest, data)
		if ok == nil {
			row[columnName] = dest
		}
		return ok
	case "DECIMAL":
		dest := new(decimal.Decimal)
		ok := ConvertAssign(&dest, data)
		if ok == nil {
			if dest != nil {
				v, _ := dest.Float64()
				row[columnName] = v
			} else {
				row[columnName] = nil
			}
		}
		return ok
	case "FLOAT", "DOUBLE":
		dest := new(float64)
		ok := ConvertAssign(&dest, data)
		if ok == nil {
			row[columnName] = dest
		}
		return ok
	case "VARCHAR", "TEXT", "CHAR":
		dest := new(string)
		ok := ConvertAssign(&dest, data)
		if ok == nil {
			row[columnName] = dest
		}
		return ok
	default:
		fmt.Printf("unhanldedtype `dbtype: %v` for column `%v` \n", columnName, dbType)
		row[columnName] = data
	}
	return nil
}
