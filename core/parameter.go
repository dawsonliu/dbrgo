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
	BodyType BodyType
	RawText  string
	Files    map[string][]*multipart.FileHeader
	BodyMap  map[string][]string
}

func (ap *ApiParams) GetMap() map[string]string {
	result := make(map[string]string)

	if ap.RequestParams != nil {
		for key, q := range ap.RequestParams {
			result[key] = strings.Join(q, ", ")
		}
	}

	if ap.BodyParams.BodyMap != nil {
		for key, q := range ap.BodyParams.BodyMap {
			result[key] = strings.Join(q, ", ")
		}
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
	if b, e := ioutil.ReadAll(c.Request.Body); e == nil {
		result.BodyParams.RawText = string(b)
		result.BodyParams.BodyType = Json

		if len(strings.TrimSpace(result.BodyParams.RawText)) > 0 {
			if json.Unmarshal(b, result.BodyParams.BodyMap) != nil {
				log.Info("failed to parse from json")
			}
		}
	} else {
		return result, e
	}

	if len(c.Request.MultipartForm.File) > 0 {
		result.BodyParams.BodyType = Files
		result.BodyParams.Files = c.Request.MultipartForm.File
	}

	if len(c.Request.Form) > 0 {
		result.BodyParams.BodyType = Forms
		result.BodyParams.BodyMap = c.Request.Form
	}

	return result, nil
}
