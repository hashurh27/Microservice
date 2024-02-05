// config/config.go
package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var dbConfig = map[string]string{
	"host":   "localhost",
	"dbname": "warranty",
	"user":   "root",
	"pass":   "",
}

// GetDBConfig returns the database configuration map
func GetDBConfig() map[string]string {
	return dbConfig
}

// CreateDBConnection creates a database connection based on the provided configuration
func CreateDBConnection(host, dbname, user, pass string) (*sql.DB, error) {
	// Create a MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, dbname)

	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// InsertData inserts data into the specified table in the database
func InsertData(db *sql.DB, tableName string, data map[string]interface{}) error {
	// Prepare SQL statement
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for col, val := range data {
		columns = append(columns, col)
		values = append(values, val)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Repeat("?, ", len(values)-1)+"?")

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}
func ReadAllData(db *sql.DB, tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range columns {
			valuePointers[i] = &values[i]
		}

		err := rows.Scan(valuePointers...)
		if err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			entry[col] = val
		}

		result = append(result, entry)
	}

	return result, nil
}

// ReadData retrieves a specific record from the specified table based on the given conditions
func ReadData(db *sql.DB, tableName string, conditions map[string]interface{}) (map[string]interface{}, error) {
	columns := make([]string, 0, len(conditions))
	values := make([]interface{}, 0, len(conditions))

	for col, val := range conditions {
		columns = append(columns, fmt.Sprintf("%s = ?", col))
		values = append(values, val)
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, strings.Join(columns, " AND "))

	rows, err := db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})

	if rows.Next() {
		columnPointers := make([]interface{}, len(columnNames))
		for i := range columnNames {
			columnPointers[i] = new(interface{})
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
		}

		for i, col := range columnNames {
			val := *(columnPointers[i].(*interface{}))
			result[col] = val
		}
	}

	return result, nil
}

// getData provides sample data
func GetData() map[string]interface{} {
	return map[string]interface{}{
		"FirstName": "John",
		"LastName":  "Doe",
		"email":     "john.doe@example.com",
		"Phone":     "1",
	}
}
