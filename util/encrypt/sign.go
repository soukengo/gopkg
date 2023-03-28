package encrypt

import (
	"sort"
	"strings"
)

func Sign(params map[string]string, key string) string {
	var keys []string
	for key := range params {
		if key != "sign" {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	var buf strings.Builder
	for i, key := range keys {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(params[key])
	}
	buf.WriteString("&")
	buf.WriteString(key)
	return strings.ToUpper(MD5(buf.String()))
}

func CheckSign(params map[string]string, key string, sign string) bool {
	localSign := Sign(params, key)
	return localSign == strings.ToUpper(sign)
}
