package main

import (
	"Warranty-Microservice/config"
	"encoding/json"
	"fmt"
	"log"
)

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
