package db2struct

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GetColumnsFromMysqlTable(dbUser string,
	dbPassword string,
	dbHost string,
	dbPort int,
	dbDatabase string,
	dbTable string) (*map[string]map[string]string, []string, error) {

	var err error
	var db *sql.DB
	if dbPassword != "" {
		db, err = sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+strconv.Itoa(dbPort)+")/"+dbDatabase+"?&parseTime=True")
	} else {
		db, err = sql.Open("mysql", dbUser+"@tcp("+dbHost+":"+strconv.Itoa(dbPort)+")/"+dbDatabase+"?&parseTime=True")
	}
	defer db.Close()

	if err != nil {
		fmt.Println("Error opening mysql database: " + err.Error())
		return nil, nil, err
	}

	columnNamesStored := []string{}
	// store column as map of maps
	columnDataTypes := make(map[string]map[string]string)
	// select column data from INFORMATION_SCHEMA
	columnDataTypeQuery := "SELECT COLUMN_NAME, COLUMN_KEY, DATA_TYPE, IS_NULLABLE, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND table_name = ? order by ordinal_position asc"

	if Debug {
		fmt.Println("running: " + columnDataTypeQuery)
	}

	rows, err := db.Query(columnDataTypeQuery, dbDatabase, dbTable)
	if err != nil {
		fmt.Println("Error selecting from db: " + err.Error())
		return nil, nil, err
	}

	if rows != nil {
		defer rows.Close()
	} else {
		return nil, nil, errors.New("No results returned for table")
	}

	for rows.Next() {
		var column string
		var columnKey string
		var dataType string
		var nullable string
		var comment string
		rows.Scan(&column, &columnKey, &dataType, &nullable, &comment)

		columnDataTypes[column] = map[string]string{"value": dataType, "nullable": nullable, "primary": columnKey, "comment": comment}
		columnNamesStored = append(columnNamesStored, column)
	}

	return &columnDataTypes, columnNamesStored, err
}

func generateMysqlTypes(obj map[string]map[string]string, columnsStored []string, dbPath int, jsonAnnotation bool, gormAnnotation bool, gureguTypes bool) string {
	structure := "struct {"
	for _, key := range columnsStored {
		mysqlType := obj[key]
		nullable := false
		if mysqlType["nullable"] == "YES" {
			nullable = true
		}

		primary := ""
		if mysqlType["primary"] == "PRI" {
			primary = ";primary_key"
		}

		// Get the corresponding go value type for this mysql type
		var valueType string
		valueType = mysqlTypeToGoType(mysqlType["value"], nullable, gureguTypes)
		fieldName := fmtFieldName(stringifyFirstChar(key))
		var annotations []string
		if gormAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", key, primary))
		}
		if jsonAnnotation == true {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", key))
		}
		if len(annotations) > 0 {
			// maybe we should add comment at the last index
			comment := mysqlType["comment"]
			structure += fmt.Sprintf("\n%s %s `%s`  //%s", fieldName, valueType, strings.Join(annotations, " "), comment)
		} else {
			structure += fmt.Sprintf("\n%s %s", fieldName, valueType)
		}
	}

	return structure
}

// mysqlTypeToGoType converts the mysql types to go compatible sql.Nullable (https://golang.org/pkg/database/sql/) types
func mysqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt
			}
			return sqlNullInt
		}
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			if gureguTypes {
				return gureguNullString
			}
			return sqlNullString
		}
		return "string"
	case "date", "datetime", "time", "timestamp":
		if nullable && gureguTypes {
			return gureguNullTime
		}
		return golangTime
	case "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat64
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat
			}
			return sqlNullFloat
		}
		return golangFloat32
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}
