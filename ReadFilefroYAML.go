package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
				fmt.Println(k, "is string", vv)
			case float64:
				fmt.Println(k, "is float64", vv)
			case []interface{}:
				fmt.Println(k, "is an array:")
				for i, u := range vv {
					fmt.Println(i+1, u)
				}
			default:
				fmt.Println(k, "is of a type I don't know how to handle")
			}
		}

	}

}
