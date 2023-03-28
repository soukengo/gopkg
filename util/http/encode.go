package http

import "strings"

func URLParams(params map[string]string) string {
	var buf strings.Builder
	i := 0
	for k, v := range params {
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
		i++
	}
	return buf.String()
}
