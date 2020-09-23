package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func PrintFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func prettyPrint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "	")
	return out.Bytes(), err
}

func main() {

	db, err := sql.Open("mysql", "root:koushi8888@tcp(127.0.0.1:3306)/student")
	PrintFatalError(err)
	defer db.Close()

	creTab, err := db.Query("CREATE TABLE IF NOT EXISTS jsons(id INT(11), author VARCHAR(500) DEFAULT NULL,description VARCHAR(500) DEFAULT NULL)")
	PrintFatalError(err)

	defer creTab.Close()

	filePath := "E:\\Go Tasks\\Target2SQL\\Target2SQL\\Test\\"

	files, err := ioutil.ReadDir(filePath)
	PrintFatalError(err)

	for _, file := range files {

		dat, err := ioutil.ReadFile(filePath + file.Name())
		PrintFatalError(err)
		fmt.Println(file.Name())
		var f interface{}

		err = json.Unmarshal([]byte(dat), &f)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		m := f.(map[string]interface{})

		for k, v := range m {
			switch vv := v.(type) {
			case string:
				if k == "author" {
					insRow, err := db.Query("INSERT INTO jsons(id,author)values(1,?)", vv)
					PrintFatalError(err)
					defer insRow.Close()

				}
				fmt.Println(k, "is string", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case []interface{}:
				fmt.Println(k, "is an array:")
				if k == "description" {
					for i, u := range vv {
						fmt.Println(i+1, u)
						insRow, err := db.Query("INSERT INTO jsons(id,description)values(1,?)", u)
						PrintFatalError(err)
						defer insRow.Close()
					}

				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}

	}

}
