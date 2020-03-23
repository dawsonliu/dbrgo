package main

import (
	"encoding/json"
	"log"
	"os"
)

type ApiModel struct {
	Version        string         `json:"version"`
	Owner          string         `json:"owner"`
	UpdatedTime    string         `json:"updatedTime"`
	Name           string         `json:"name"`
	Module         string         `json:"module"`
	URL            string         `json:"url"`
	UseAbsoluteURL bool           `json:"useAbsoluteUrl"`
	Method         string         `json:"method"`
	Title          string         `json:"title"`
	Summary        string         `json:"summary"`
	Note           string         `json:"note"`
	AllowAnonymous bool           `json:"allowAnonymous"`
	Cache          Cache          `json:"cache"`
	Implemented    bool           `json:"implemented"`
	Implementation Implementation `json:"implementation"`
	Parameter      Parameter      `json:"parameter"`
	Result         Result         `json:"result"`
	Mock           []Mock         `json:"mock"`
}
type Cache struct {
	Enabled    bool   `json:"enabled"`
	Type       string `json:"type"`
	Expiration int    `json:"expiration"`
}
type Implementation struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
type Body struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
type Parameter struct {
	Query []interface{} `json:"query"`
	Body  []Body        `json:"body"`
}
type Schema struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Nullable    bool     `json:"nullable"`
	Description string   `json:"description"`
	Schema      []Schema `json:"schema,omitempty"`
}
type Result struct {
	Type   string   `json:"type"`
	Schema []Schema `json:"schema"`
}
type Mock struct {
	Input  map[string]string      `json:"input"`
	Output map[string]interface{} `json:"output"`
}

func LoadApis() {
	cd, _ := os.Getwd()
	file, err := os.Open(cd + "..\\..\\resources\\api\\api.json")
	if err != nil {
		log.Fatal("opening config file", err.Error())
	}

	defer file.Close()

	var apis []ApiModel

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&apis); err != nil {
		log.Fatal("parsing config file", err.Error())
	}

	log.Println(apis)
}
