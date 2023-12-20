package utils

import (
	"strconv"
	"strings"
)

// unicodeToText
// @description: unicode转文本
// @param str
// @return dst
func unicodeToText(str string) (dst string) {
	dst, _ = strconv.Unquote(strings.Replace(strconv.Quote(str), `\\u`, `\u`, -1))
	return
}
