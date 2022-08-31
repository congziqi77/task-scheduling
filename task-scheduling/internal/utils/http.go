package utils

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func Get(uri string, route *gin.Engine) *httptest.ResponseRecorder {

	//构造get请求
	req := httptest.NewRequest("GET", uri, nil)

	//初始化相应
	w := httptest.NewRecorder()

	//调用相应的handler接口
	route.ServeHTTP(w, req)

	return w
}

func PostJson(uri string, param map[string]interface{}, route *gin.Engine) *httptest.ResponseRecorder {
	jsonByte, err := json.Marshal(param)
	if err != nil {
		return nil
	}

	req := httptest.NewRequest("POST", uri, bytes.NewBuffer(jsonByte))
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	return w
}
