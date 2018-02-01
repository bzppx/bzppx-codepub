package controllers

import (
	"bzppx-codepub/app/models"
	"bzppx-codepub/app/utils"
	"regexp"
	"strings"
	"time"
)

type ApiAuthController struct {
	BaseController
}

func (this *ApiAuthController) List() {
	page, _ := this.GetInt("page", 1)

	number := 20
	limit := (page - 1) * number
	var err error
	var count int64
	var authApis []map[string]string

	authApis, err = models.ApiAuthModel.GetApiAuthByLimit(number, limit)
	count, err = models.ApiAuthModel.CountApiAuth()
	if err != nil {
		this.viewError(err.Error(), "api/list")
	}

	this.Data["authApis"] = authApis
	this.SetPaginator(number, count)

	this.viewLayoutTitle("auth_api列表", "api/list", "page")
}

func (this *ApiAuthController) Add() {
	this.viewLayoutTitle("添加auth_api", "api/form", "page")
}

func (this *ApiAuthController) Edit() {
	authApiId := strings.Trim(this.GetString("api_auth_id", ""), "")
	authApi, err := models.ApiAuthModel.GetApiAuthByApiAuthId(authApiId)
	if err != nil {
		this.viewError(err.Error(), "api/form")
	}
	if len(authApi) == 0 {
		this.viewError("未查到api信息", "api/form")
	}

	this.Data["apiAuth"] = authApi
	this.viewLayoutTitle("修改auth_api", "api/form", "page")
}

func (this *ApiAuthController) Save() {
	name := strings.Trim(this.GetString("name", ""), "")
	url := strings.Trim(this.GetString("url", ""), "")
	key := strings.Trim(this.GetString("key", ""), "")
	sort := strings.Trim(this.GetString("sort", ""), "")
	isShow := strings.Trim(this.GetString("is_show", ""), "")

	if name == "" {
		this.jsonError("api名称不能为空！")
	}
	ok, err := regexp.MatchString(`(https?)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, url)
	if err != nil {
		this.jsonError("url正则匹配失败！")
	}
	if !ok {
		this.jsonError("url不正确！")
	}
	if key == "" {
		this.jsonError("key不能为空！")
	}
	exist, err := models.ApiAuthModel.CheckKeyExist(key, "")
	if exist {
		this.jsonError("key已存在！")
	}
	if sort == "" {
		this.jsonError("排序号不能为空！")
	}

	timeNow := time.Now().Unix()
	apiAuth := map[string]interface{}{
		"name":        name,
		"url":         url,
		"key":         key,
		"sort":        sort,
		"is_show":     isShow,
		"is_delete":   0,
		"create_time": timeNow,
		"update_time": timeNow,
	}
	authApiId, err := models.ApiAuthModel.Insert(apiAuth)
	if err != nil {
		this.ErrorLog("添加auth_api失败: " + err.Error())
		this.jsonError("添加auth_api失败！")
	}

	this.InfoLog("添加auth_api " + utils.NewConvert().IntToString(authApiId, 10) + " 成功")
	this.jsonSuccess("添加auth_api成功", nil, "/apiAuth/list")
}

// 修改 api_auth
func (this *ApiAuthController) Modify() {
	authApiId := strings.Trim(this.GetString("api_auth_id", ""), "")
	name := strings.Trim(this.GetString("name", ""), "")
	url := strings.Trim(this.GetString("url", ""), "")
	sort := strings.Trim(this.GetString("sort", ""), "")
	isShow := strings.Trim(this.GetString("is_show", ""), "")

	if name == "" {
		this.jsonError("api名称不能为空！")
	}
	if authApiId == "" {
		this.jsonError("参数错误！")
	}
	ok, err := regexp.MatchString(`(https?)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`, url)
	if err != nil {
		this.jsonError("url正则匹配失败！")
	}
	if !ok {
		this.jsonError("url不正确！")
	}
	if sort == "" {
		this.jsonError("排序号不能为空！")
	}

	timeNow := time.Now().Unix()
	apiAuth := map[string]interface{}{
		"name":        name,
		"url":         url,
		"sort":        sort,
		"is_show":     isShow,
		"is_delete":   0,
		"update_time": timeNow,
	}

	_, err = models.ApiAuthModel.UpdateByAuthApiId(apiAuth, authApiId)
	if err != nil {
		this.ErrorLog("修改auth_api失败: " + err.Error())
		this.jsonError("修改auth_api失败！")
	}
	this.InfoLog("修改auth_api " + authApiId + " 成功")
	this.jsonSuccess("修改auth_api成功", nil, "/apiAuth/list")
}

// 删除 api_auth
func (this *ApiAuthController) Delete() {
	authApiId := strings.Trim(this.GetString("api_auth_id", ""), "")

	if authApiId == "" {
		this.jsonError("参数错误")
	}
	_, err := models.ApiAuthModel.DeleteByAuthApiId(authApiId)
	if err != nil {
		this.ErrorLog("删除auth_api失败: " + err.Error())
		this.jsonError("删除auth_api失败！")
	}

	this.InfoLog("删除auth_api成功")
	this.jsonSuccess("删除auth_api成功", nil, "/apiAuth/list")
}

// api 接入手册
func (this *ApiAuthController) Introduce() {
	this.viewLayoutTitle("api接入手册", "api/introduce", "page")
}
