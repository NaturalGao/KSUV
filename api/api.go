package api

import (
	"bytes"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type ResultMp map[string]interface{}

const API_HOST = "https://cp.kuaishou.com/rest/cp/works/v2/video/pc/"

var ClientHttp *resty.Client

var ApiLists map[string]string = map[string]string{
	"relationList": "relation/list",
	"uploadFinish": "upload/finish",
	"submit":       "submit",
}

type ApiObject struct {
	Cookie string
}

func New(Cookie string) *ApiObject {
	a := new(ApiObject)
	a.Cookie = Cookie
	return a
}

func init() {
	ClientHttp = resty.New()
}

func (a *ApiObject) RelationList(body map[string]interface{}) (ResultMp, error) {
	resp, err := ClientHttp.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", a.Cookie).
		SetBody(body).
		Post(a.GetRequestUrl("relationList"))

	if err != nil {
		return make(map[string]interface{}), err
	}

	r, err := JsonToMap(resp.Body())

	return r, err
}

/**
获取视频 信息
*/
func (a *ApiObject) UploadFinish(body map[string]interface{}) (ResultMp, error) {
	resp, err := ClientHttp.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", a.Cookie).
		SetBody(body).
		Post(a.GetRequestUrl("uploadFinish"))

	if err != nil {
		return make(map[string]interface{}), err
	}

	r, err := JsonToMap(resp.Body())

	return r, err
}

// 发布视频
func (a *ApiObject) SubmitVideo(body map[string]interface{}) (ResultMp, error) {
	resp, err := ClientHttp.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Cookie", a.Cookie).
		SetBody(body).
		Post(a.GetRequestUrl("submit"))

	if err != nil {
		return make(map[string]interface{}), err
	}

	r, err := JsonToMap(resp.Body())

	return r, err
}

/**
上传视频
*/
func (a *ApiObject) UploadMultipart(up_token string, fileName string, fileBytes []byte) (ResultMp, error) {
	resp, err := ClientHttp.R().
		SetHeader("Cookie", a.Cookie).
		SetQueryParam("upload_token", up_token).
		SetFileReader("file", fileName, bytes.NewReader(fileBytes)).
		Post("https://js-cp-upload.xxpkg.com/api/upload/multipart")

	if err != nil {
		return make(map[string]interface{}), err
	}

	r, err := JsonToMap(resp.Body())

	return r, err
}

func (a *ApiObject) GetRequestUrl(key string) string {
	return API_HOST + ApiLists[key]
}

func JsonToMap(body []byte) (tempResultMap ResultMp, err error) {

	err = json.Unmarshal(body, &tempResultMap)

	if err != nil {
		return tempResultMap, err
	}

	return tempResultMap, nil
}
