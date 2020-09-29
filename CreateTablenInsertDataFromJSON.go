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
		s := make([]string, len(vv))
		for i, u := range vv {
			s[i] = fmt.Sprint(u)
		}
		str := strings.Join(s, "**")
		return str
	case interface{}:
		str := fmt.Sprintf("%v", t)
		return str
	default:
		return "null"
	}
}

func main() {
	count := 0
	db, err := sql.Open("mysql", "root:xxxxx@tcp(127.0.0.1:3306)/student")
	PrintFatalError(err)
	defer db.Close()

	modTab, err := db.Query("CREATE TABLE IF NOT EXISTS modulesinfo( module_name VARCHAR(500) DEFAULT NULL, module_version_added VARCHAR(50), module_author VARCHAR(500) DEFAULT NULL, module_short_description VARCHAR(500) DEFAULT NULL, module_extends_documentation_fragment VARCHAR(500) DEFAULT NULL,  module_requirements VARCHAR(500) DEFAULT NULL, module_description VARCHAR(500) DEFAULT NULL, module_deprecated VARCHAR(500) DEFAULT NULL, option_name VARCHAR(500), option_version_added VARCHAR(50) DEFAULT NULL, option_description VARCHAR(5000), option_required VARCHAR(100) DEFAULT NULL, option_defaults VARCHAR(1000) DEFAULT NULL, option_choices VARCHAR(100) DEFAULT NULL, option_suboptions VARCHAR(1000) DEFAULT NULL, option_type VARCHAR(1000) DEFAULT NULL,  option_elements VARCHAR(1000) DEFAULT NULL, option_aliases VARCHAR(100)DEFAULT NULL )ENGINE=InnoDB DEFAULT CHARSET=utf8")
	PrintFatalError(err)
	defer modTab.Close()

	filePath := "E:\\Go Tasks\\Final\\Results\\JSONOut\\"
	files, err := ioutil.ReadDir(filePath)
	PrintFatalError(err)
	TotalCount := len(files)
	fmt.Println("Number of files in the Directory: ", TotalCount)

	for _, file := range files {
		count++
		raw, err := ioutil.ReadFile(filePath + file.Name())
		PrintFatalError(err)

		var initial Initial
		var options map[string]Options
		//fmt.Println("File Name with Path :-->> ", filePath+file.Name())
		if err := json.Unmarshal([]byte(raw), &initial); err != nil {
			log.Fatal(err)
		}
		b, err := json.Marshal(initial.Options)
		PrintFatalError(err)
		if err := json.Unmarshal([]byte(b), &options); err != nil {
			log.Fatal(err)
		}

		keys := reflect.ValueOf(options).MapKeys()
		for i := 0; i < len(keys); i++ {
			OptName := keys[i].String()
			OptVersionAdded := options[keys[i].String()].VersionAdded
			OptRequired := options[keys[i].String()].Required
			OptElements := options[keys[i].String()].Elements
			OptType := options[keys[i].String()].Type
			OptAlias := strings.Join(options[keys[i].String()].Aliases, "->>")
			OptDesc := getValueswithin(options[keys[i].String()].Description)
			OptDef := getValueswithin(options[keys[i].String()].Default)
			OptChoice := getValueswithin(options[keys[i].String()].Choices)
			OptSuboptions := getValueswithin(options[keys[i].String()].Suboptions)
			IniAuthor := getValueswithin(initial.Author)
			IniReq := getValueswithin(initial.Requirements)
			IniDesc := strings.Join(initial.Description, "->>")
			IniExtDoc := strings.Join(initial.ExtendsDoc, "->>")

			iniForm, err := db.Prepare("INSERT INTO modulesinfo (module_name, module_version_added, module_author, module_short_description, module_extends_documentation_fragment, module_requirements, module_description, module_deprecated, option_name, option_version_added, option_description, option_required, option_defaults, option_choices, option_suboptions, option_type, option_elements, option_aliases) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
			PrintFatalError(err)

			iniForm.Exec(initial.Module, initial.VersionAdded, IniAuthor, initial.ShortDescription, IniExtDoc, IniReq, IniDesc, initial.Deprecated, OptName, OptVersionAdded, OptDesc, OptRequired, OptDef, OptChoice, OptSuboptions, OptType, OptElements, OptAlias)
			defer iniForm.Close()
		}
		fmt.Printf("Inserted File: %d Out of: %d\n", count, TotalCount)
	}

}
