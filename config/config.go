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
	"pass":   "", // Consider loading this from a secure source in a production environment.
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
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set up connection pooling
	db.SetMaxOpenConns(10) // Adjust as needed
	db.SetMaxIdleConns(5)  // Adjust as needed

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// ConnectedDb creates a database connection using the configuration in dbConfig
func ConnectedDb() (*sql.DB, error) {
	return CreateDBConnection(dbConfig["host"], dbConfig["dbname"], dbConfig["user"], dbConfig["pass"])
}

// InsertData inserts data into the specified table in the database
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

// getUserData provides sample data
func GetUserData() map[string]interface{} {
	return map[string]interface{}{
		"FirstName": "John",
		"LastName":  "Doe",
		"email":     "john.doe@example.com",
		"Phone":     "1",
	}
}

/*
usage

func main() {
	// Get the database configuration map
	dbConfig := config.GetDBConfig()

	// Create a new database connection
	db, err := config.CreateDBConnection(dbConfig["host"], dbConfig["dbname"], dbConfig["user"], dbConfig["pass"])
	if err != nil {
		log.Fatal(err)
	}

	// Insert data into the database
	data := config.GetData()
	tableName := "customer"

	err = config.InsertData(db, tableName, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Insert successful.")

	// Read all data
	allData, err := config.ReadAllData(db, tableName)
	if err != nil {
		log.Fatal(err)
	}
	printJSON("All data:", allData)

	// Read specific data
	conditions := map[string]interface{}{
		"FirstName": "John",
	}
	specificData, err := config.ReadData(db, tableName, conditions)
	if err != nil {
		log.Fatal(err)
	}
	printJSON("Specific data:", specificData)
}

// printJSON prints the data in JSON format
func printJSON(message string, data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:\n%s\n", message, jsonData)
}

*/
