package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

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

	var author string
	var shortDescription string
	var versionAdded string
	var extendsDoc string
	var module string
	var description string

	db, err := sql.Open("mysql", "root:koushi8888@tcp(127.0.0.1:3306)/student")
	PrintFatalError(err)
	defer db.Close()

	creTab, err := db.Query("CREATE TABLE IF NOT EXISTS jsons(module VARCHAR(500), short_description VARCHAR(500), version_added VARCHAR(500), extends_document VARCHAR(1000), author VARCHAR(500) DEFAULT NULL,description VARCHAR(4000) DEFAULT NULL)")
	PrintFatalError(err)

	defer creTab.Close()

	//filePath := "E:\\Go Tasks\\Target2SQL\\Target2SQL\\Test\\"

	filePath := "E:\\Go Tasks\\Final\\Results\\JSONOut\\"

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
					author = vv
				} else if k == "module" {
					module = vv
				} else if k == "short_description" {
					shortDescription = vv
				} else if k == "version_added" {
					versionAdded = vv
				}
				fmt.Println(k, "is string", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case []interface{}:
				fmt.Println(k, "is an array:")
				if k == "description" {
					s := make([]string, len(vv))
					for i, u := range vv {
						fmt.Println(i+1, u)
						s[i] = fmt.Sprint(u)
					}
					description = strings.Join(s, "->>")
				} else if k == "extends_documentation_fragment" {
					s := make([]string, len(vv))
					for i, u := range vv {
						fmt.Println(i+1, u)
						s[i] = fmt.Sprint(u)
					}
					extendsDoc = strings.Join(s, "->>")
				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}
		insRow, err := db.Query("INSERT INTO jsons(module,short_description,version_added,extends_document,author,description)values(?,?,?,?,?,?)", module, shortDescription, versionAdded, extendsDoc, author, description)
		PrintFatalError(err)
		defer insRow.Close()
	}

}
