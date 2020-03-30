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
	BodyType    BodyType // Body中数据类型
	rawJsonText string   // 如果是Json类型，则存储json字符串，需要的时候可以通过 `Json()` 和 `JsonTo(interface{})` 获取
	Files       map[string][]*multipart.FileHeader
	BodyMap     interface{}
}

func (ap *ApiParams) Json() (interface{}, error) {
	result := new(interface{})

	if err := json.Unmarshal([]byte(ap.BodyParams.rawJsonText), &result); err != nil {
		log.Info("failed to parse from json" + err.Error())
		return nil, err
	}

	return result, nil
}

// 将Json序列化成指定的强类型 (Struct)
func (ap *ApiParams) JsonTo(data interface{}) error {
	if err := json.Unmarshal([]byte(ap.BodyParams.rawJsonText), &data); err != nil {
		log.Info("failed to parse from json" + err.Error())
		return err
	}

	return nil
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
		result.BodyParams.rawJsonText = string(b)
		result.BodyParams.BodyType = Json
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
