package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// ConvertReq 转换请求
type ConvertReq struct {
	Ddl         string `json:"ddl"`
	DdlType     string `json:"ddlType"`     // mysql ddl
	ConvertType string `json:"convertType"` // dm
}

// Resp 返回实体
type Resp struct {
	Code string `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

// Convert ddl convert
func Convert(writer http.ResponseWriter, request *http.Request) {
	var resp Resp
	var req ConvertReq

	// json -> ConvertReq
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		log.Printf("json转换异常：{%v}", err)
		resp.Msg = "请求json异常"
		resp.Code = "500"
		resp.response(writer)
		return
	}

	// check DDL
	if len(req.Ddl) == 0 {
		resp.response(writer)
		return
	}

	resp.Code = "200"
	resp.Data = mysqlToDm(req.Ddl)
	resp.Msg = "success"
	resp.response(writer)
}

// response data
func (resp Resp) response(writer http.ResponseWriter) {
	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		return
	}
}

// handle rule
func mapRule() map[string]string {
	// rule
	return map[string]string{
		"`":                               "\"",
		"(?i)longtext":                    "CLOB",
		"(?i)ON UPDATE CURRENT_TIMESTAMP": "",
		"(?i)NOT NULL":                    "",
		"(?i),?\\s*\\n*PRIMARY\\s+KEY\\s*\\([^)]+\\)": "",
		"(?i)(int)\\s*\\(\\s*\\d+\\s*\\)":             "INT",
		"(?i)(bigint)\\s*\\(\\s*\\d+\\s*\\)":          "BIGINT",
		"(?i)(tinyint)\\s*\\(\\s*\\d+\\s*\\)":         "TINYINT",
	}
}

// mysqlToDm
func mysqlToDm(ddl string) string {
	// rule
	for oldKey, newKey := range mapRule() {
		regex := regexp.MustCompile(oldKey)
		ddl = regex.ReplaceAllString(ddl, newKey)
	}

	// remove PRIMARY KEY (`id`)
	startIndex := strings.Index(ddl, "CREATE")
	if startIndex < 0 {
		startIndex = strings.Index(ddl, "create")
	}
	endIndex := strings.LastIndex(ddl, ")")
	ddl = ddl[startIndex : endIndex+1]
	ddl = strings.ReplaceAll(ddl, "AUTO_INCREMENT", "PRIMARY KEY IDENTITY(1, 1) ")
	return ddl
}
