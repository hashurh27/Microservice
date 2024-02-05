package forms

import (
	"Warranty-Microservice/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CarSender(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hellow world")
}

func CustomerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		CarSender(w, r)

	case http.MethodPost:
		var data map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		jsonData, err := json.Marshal(data)

		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Print the JSON data as a string
		fmt.Println("Re-encoded JSON:", string(jsonData))
		fmt.Println("Re-encoded JSON:", data)

		dbConfig := config.GetDBConfig()

		// Create a new database connection
		db, err := config.CreateDBConnection(dbConfig["host"], dbConfig["dbname"], dbConfig["user"], dbConfig["pass"])
		if err != nil {
			log.Fatal(err)
		}

		// Insert data into the database
		//data := config.GetData()
		tableName := "customer"

		err = config.InsertData(db, tableName, data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Insert successful.")
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
