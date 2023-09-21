package api

import (
	"encoding/json"
	"log"
	"net/http"
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

	// dm rule handle
	req.Ddl = strings.ReplaceAll(req.Ddl, "`", "\"")
	req.Ddl = strings.ReplaceAll(req.Ddl, "longtext", "clob")
	req.Ddl = strings.ReplaceAll(req.Ddl, "ON UPDATE CURRENT_TIMESTAMP", "")
	req.Ddl = strings.ReplaceAll(req.Ddl, "AUTO_INCREMENT", "IDENTITY (1, 1)")
	startIndex := strings.Index(req.Ddl, "CREATE")
	if startIndex < 0 {
		startIndex = strings.Index(req.Ddl, "create")
	}
	endIndex := strings.LastIndex(req.Ddl, ")")
	dm := req.Ddl[startIndex : endIndex+1]

	resp.Code = "200"
	resp.Data = dm
	resp.Msg = "success"
	resp.response(writer)
}

// response data
func (resp Resp) response(writer http.ResponseWriter) {
	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		return
	}
}
