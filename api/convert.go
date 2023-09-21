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
		"`":                           "\"",
		"longtext":                    "clob",
		"ON UPDATE CURRENT_TIMESTAMP": "",
		"AUTO_INCREMENT":              "IDENTITY(1, 1) PRIMARY KEY",
	}
}

// mysqlToDm
func mysqlToDm(ddl string) string {
	// rule
	for oldKey, newKey := range mapRule() {
		ddl = strings.ReplaceAll(ddl, oldKey, newKey)
		log.Printf("规则转换结果：{%s}", ddl)
	}

	log.Print(ddl)
	// remove PRIMARY KEY (`id`)
	startIndex := strings.Index(ddl, "CREATE")
	if startIndex < 0 {
		startIndex = strings.Index(ddl, "create")
	}
	endIndex := strings.LastIndex(ddl, ")")
	ddl = ddl[startIndex : endIndex+1]
	arr := strings.Split(ddl, ",")
	var dm strings.Builder
	for _, str := range arr {
		regex := regexp.MustCompile(`PRIMARY\s+KEY\s*\([^)]+\)`)
		if !regex.MatchString(str) {
			dm.WriteString(str)
		}
	}
	return dm.String()
}
