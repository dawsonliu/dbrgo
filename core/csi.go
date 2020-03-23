package main

import (
	"encoding/json"
	"log"
	"os"
)

type CsiModal struct {
	Name                string      `json:"name"`
	Code                Code        `json:"code"`
	Module              string      `json:"module"`
	ResultSet           string      `json:"resultSet"`
	QueryOnly           bool        `json:"queryOnly"`
	RequiredTransaction bool        `json:"requiredTransaction"`
	Middlewares         Middlewares `json:"middlewares"`
}
type Pagination struct {
	Size  string `json:"size"`
	Count string `json:"count"`
	Page  string `json:"page"`
}
type Middlewares struct {
	Pagination Pagination `json:"pagination"`
}
type Code struct {
	IsObject  bool
	CodeMap   map[string]string
	PlainCode string
}

func LoadCsis() error {
	cd, _ := os.Getwd()
	file, err := os.Open(cd + "..\\..\\resources\\csi\\csi.json")
	if err != nil {
		log.Fatal("failed to open csi file", err.Error())
		return err
	}

	defer file.Close()

	var csis []CsiModal

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&csis); err != nil {
		log.Fatal("failed to parse csi file", err.Error())
		return err
	}

	return nil
}

func (code Code) MarshalJSON() ([]byte, error) {
	if code.IsObject {
		return json.Marshal(code.CodeMap)
	} else {
		return []byte(code.PlainCode), nil
	}
}

func (code *Code) UnmarshalJSON(data []byte) error {
	var cm map[string]string
	if err := json.Unmarshal(data, &cm); err != nil {
		*code = Code{IsObject: false, PlainCode: string(data)}
	} else {
		*code = Code{IsObject: true, CodeMap: cm}
	}

	return nil
}
