package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func readConnectionString(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading connection file: %v", err)
	}
	return strings.TrimSpace(string(content)), nil
}

func main() {
	connStr, err := readConnectionString(".qex")
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	}
	reader := bufio.NewReader(os.Stdin)
	query, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error reading query: ", err)
	}
	query = strings.TrimSpace(query)
	variables := input(query)
	fmt.Print("\n")
	rows, err := db.Query(query, variables...)
	if err != nil {
		log.Fatal("Error executing query: ", err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal("Error getting column names: ", err)
	}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	var result []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal("Error scanning row: ", err)
		}
		row := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			switch val := val.(type) {
			case nil:
				v = nil
			case []byte:
				v = string(val)
			case time.Time:
				v = val.Format("2006-01-02 15:04:05")
			default:
				v = val
			}
			row[col] = v
		}
		result = append(result, row)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("Error iterating rows: ", err)
	}
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal("Error converting to JSON: ", err)
	}
	fmt.Println(string(jsonData))
}
