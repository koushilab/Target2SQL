package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func PrintFatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Initial struct {
	Module           string      `json:"module,omitempty"`
	Options          interface{} `json:"options,omitempty"`
	ShortDescription string      `json:"short_description,omitempty"`
	Description      []string    `json:"description,omitempty"`
	ExtendsDoc       []string    `json:"extends_documentation_fragment,omitempty"`
	VersionAdded     string      `json:"version_added,omitempty"`
	Author           interface{} `json:"author,omitempty"`
	Deprecated       string      `json:"deprecated,omitempty"`
	Requirements     interface{} `json:"requirements,omitempty"`
}

type Options struct {
	Description  interface{} `json:"description,omitempty"`
	Required     bool        `json:"required,omitempty"`
	Default      interface{} `json:"default,omitempty"`
	Choices      interface{} `json:"choices,omitempty"`
	Type         string      `json:"type,omitempty"`
	Elements     string      `json:"elements,omitempty"`
	Aliases      []string    `json:"aliases,omitempty"`
	VersionAdded string      `json:"version_added,omitempty"`
	Suboptions   interface{} `json:"suboptions,omitempty"`
}

func getValueswithin(t interface{}) string {

	switch vv := t.(type) {
	case []interface{}:
		//	fmt.Println(t, " is an array:")
		s := make([]string, len(vv))
		for i, u := range vv {
			//fmt.Println(i+1, u)
			s[i] = fmt.Sprint(u)
		}
		str := strings.Join(s, "**")
		//	fmt.Println(description)
		return str
	case interface{}:
		str := fmt.Sprintf("%v", t)
		//	fmt.Println(reflect.ValueOf(t))
		return str
	default:
		return "null"
	}
}

func main() {

	count := 0
	db, err := sql.Open("mysql", "root:xxxxxx@tcp(127.0.0.1:3306)/student")
	PrintFatalError(err)
	defer db.Close()

	optTab, err := db.Query("CREATE TABLE IF NOT EXISTS options(module VARCHAR(500),name VARCHAR(500), version_added VARCHAR(50), description VARCHAR(5000), required VARCHAR(100), defaults VARCHAR(1000), choices VARCHAR(100), elements VARCHAR(1000), aliases VARCHAR(100), suboptions VARCHAR(1000), type VARCHAR(1000) DEFAULT NULL)")
	PrintFatalError(err)
	defer optTab.Close()

	modTab, err := db.Query("CREATE TABLE IF NOT EXISTS modules( module VARCHAR(500) DEFAULT NULL, version_added VARCHAR(50), author VARCHAR(500) DEFAULT NULL, short_description VARCHAR(500) DEFAULT NULL, extends_documentation_fragment VARCHAR(500) DEFAULT NULL,  requirements VARCHAR(500) DEFAULT NULL, description VARCHAR(500) DEFAULT NULL, deprecated VARCHAR(500) DEFAULT NULL)ENGINE=InnoDB DEFAULT CHARSET=utf8")
	PrintFatalError(err)
	defer modTab.Close()

	filePath := "E:\\Go Tasks\\Final\\Results\\JSONOut\\"

	files, err := ioutil.ReadDir(filePath)
	PrintFatalError(err)

	for _, file := range files {
		count++
		raw, err := ioutil.ReadFile(filePath + file.Name())
		PrintFatalError(err)

		var initial Initial

		var options map[string]Options
		fmt.Println("File Name with Path :-->> ", filePath+file.Name())

		if err := json.Unmarshal([]byte(raw), &initial); err != nil {
			log.Fatal(err)
		}

		b, err := json.Marshal(initial.Options)
		if err != nil {
			fmt.Println("error:", err)
		}

		if err := json.Unmarshal([]byte(b), &options); err != nil {
			log.Fatal(err)
		}
		// Convert string[] to string joins
		//  ******** InitialDesc := strings.Join(initial.Description, "->>")
		//  ******** InitialExtendsDoc := strings.Join(initial.ExtendsDoc, "->>")

		keys := reflect.ValueOf(options).MapKeys()
		fmt.Println("Module: ", initial.Module)
		for i := 0; i < len(keys); i++ {
			s := keys[i].String()
			OptVersionAdded := options[keys[i].String()].VersionAdded
			OptRequired := options[keys[i].String()].Required
			OptElements := options[keys[i].String()].Elements
			OptType := options[keys[i].String()].Type
			OptAlias := strings.Join(options[keys[i].String()].Aliases, "->>")
			OptDesc := getValueswithin(options[keys[i].String()].Description)
			OptDef := getValueswithin(options[keys[i].String()].Default)
			OptChoice := getValueswithin(options[keys[i].String()].Choices)
			OptSuboptions := getValueswithin(options[keys[i].String()].Suboptions)

			insForm, err := db.Prepare("INSERT INTO options (module, name, version_added, description, required, defaults, choices, elements, aliases, suboptions, type) VALUES(?,?,?,?,?,?,?,?,?,?,?)")
			PrintFatalError(err)
			insForm.Exec(initial.Module, s, OptVersionAdded, OptDesc, OptRequired, OptDef, OptChoice, OptElements, OptAlias, OptSuboptions, OptType)
			defer insForm.Close()

		}

		IniAuthor := getValueswithin(initial.Author)
		IniReq := getValueswithin(initial.Requirements)
		IniDesc := strings.Join(initial.Description, "->>")
		IniExtDoc := strings.Join(initial.ExtendsDoc, "->>")

		iniForm, err := db.Prepare("INSERT INTO modules (module, version_added, author, short_description, extends_documentation_fragment, requirements, description, deprecated) VALUES(?,?,?,?,?,?,?,?)")
		PrintFatalError(err)

		iniForm.Exec(initial.Module, initial.VersionAdded, IniAuthor, initial.ShortDescription, IniExtDoc, IniReq, IniDesc, initial.Deprecated)
		defer iniForm.Close()
		fmt.Println(count)
	}
}
