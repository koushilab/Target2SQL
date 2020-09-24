package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

type Initial struct {
	Module           string      `json:"module,omitempty"`
	Options          interface{} `json:"options,omitempty"`
	ShortDescription string      `json:"short_description,omitempty"`
	Description      []string    `json:"description,omitempty"`
	VersionAdded     string      `json:"version_added,omitempty"`
	Author           string      `json:"author,omitempty"`
	Deprecated       string      `json:"deprecated,omitempty"`
	Requirements     string      `json:"requirements,omitempty"`
}

type Options struct {
	Description  []string `json:"description,omitempty"`
	Required     bool     `json:"required,omitempty"`
	Default      string   `json:"default,omitempty"`
	Choices      []string `json:"choices,omitempty"`
	Type         string   `json:"type,omitempty"`
	Elements     string   `json:"elements,omitempty"`
	Aliases      []string `json:"aliases,omitempty"`
	VersionAdded string   `json:"version_added,omitempty"`
	Suboptions   string   `json:"suboptions,omitempty"`
}

func main() {
	raw, err := ioutil.ReadFile("E:\\Go Tasks\\Target2SQL\\Target2SQL\\Test\\optionsEx4.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	var initial Initial

	var options map[string]Options

	if err := json.Unmarshal([]byte(raw), &initial); err != nil {
		log.Fatal(err)
	}

	//	fmt.Println(initial.Options)

	b, err := json.Marshal(initial.Options)
	if err != nil {
		fmt.Println("error:", err)
	}
	//os.Stdout.Write(b)

	if err := json.Unmarshal([]byte(b), &options); err != nil {
		log.Fatal(err)
	}

	keys := reflect.ValueOf(options).MapKeys()
	fmt.Println(keys)
	fmt.Println(len(keys))
	fmt.Println(options[keys[6].String()].Description)
	fmt.Println(initial.Module)
	fmt.Println(initial.Author)
	fmt.Println(initial.Deprecated)
	//fmt.Printf("%+v\n", options)

}
