// Package utils
// @Title        tool.go
// @Description
// @Author       gxk
// @Time         2023/11/27 7:07 PM
package utils

import "encoding/json"

func Json(data interface{}) string {
	s, _ := json.Marshal(data)
	return string(s)
}
