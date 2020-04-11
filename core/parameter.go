package core

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type BodyType int

const (
	_              = iota
	Json  BodyType = iota // Json text
	Forms                 // Data from `multipart/form-data` or `application/x-www-form-urlencoded`
	Files                 // Files
	Text                  // PlainText
)

type ApiParams struct {
	RequestParams map[string][]string
	BodyParams    BodyParams
}

type BodyParams struct {
	BodyType    BodyType
	RawJsonText string
	Files       map[string][]*multipart.FileHeader
	BodyMap     interface{}
}

func (ap *ApiParams) GetMap() map[string]string {
	result := make(map[string]string)

	if ap.RequestParams != nil {
		for key, q := range ap.RequestParams {
			result[key] = strings.Join(q, ", ")
		}
	}

	if ap.BodyParams.BodyMap != nil {
		switch ap.BodyParams.BodyMap.(type) {
		case map[string][]string:
			mapData := ap.BodyParams.BodyMap.(map[string][]string)
			for key, q := range mapData {
				result[key] = strings.Join(q, ", ")
			}
		}
	}

	return result
}

func (ap *ApiParams) GetJson() interface{} {
	if len(ap.BodyParams.RawJsonText) <= 0 {
		return nil
	}

	var result interface{}
	if err := json.Unmarshal([]byte(ap.BodyParams.RawJsonText), &result); err != nil {
		log.Info("failed to parse from json" + err.Error())
	}

	return result
}

func (ap *ApiParams) Parse(dest *interface{}) error {
	param := ap.GetMap()
	return mapstructure.Decode(param, dest)
}

func ExtractParams(c *gin.Context) (ApiParams, error) {
	result := ApiParams{}

	if len(c.Request.URL.Query()) > 0 {
		result.RequestParams = c.Request.URL.Query()
	}

	result.BodyParams = BodyParams{}
	if b, e := ioutil.ReadAll(c.Request.Body); e == nil && len(b) > 0 {
		result.BodyParams.RawJsonText = string(b)
		result.BodyParams.BodyType = Json

		if len(strings.TrimSpace(result.BodyParams.RawJsonText)) > 0 {
			if err := json.Unmarshal(b, &result.BodyParams.BodyMap); err != nil {
				log.Info("failed to parse from json" + err.Error())
			}
		}
	}

	if c.Request.MultipartForm != nil && len(c.Request.MultipartForm.File) > 0 {
		result.BodyParams.BodyType = Files
		result.BodyParams.Files = c.Request.MultipartForm.File
	}

	if len(c.Request.Form) > 0 {
		result.BodyParams.BodyType = Forms
		result.BodyParams.BodyMap = c.Request.Form
	}

	return result, nil
}
