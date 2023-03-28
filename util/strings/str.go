package strings

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

//DeleteExtraSpace 删除字符串中的多余空格，有多个空格时，仅保留一个空格
func DeleteExtraSpace(s string) string {
	s1 := strings.ReplaceAll(s, "\t", " ")
	regStr := "\\s{2,}"
	reg, _ := regexp.Compile(regStr)
	s2 := make([]byte, len(s1))
	copy(s2, s1)
	spcIndex := reg.FindStringIndex(string(s2))
	for len(spcIndex) > 0 { //找到适配项
		s2 = append(s2[:spcIndex[0]+1], s2[spcIndex[1]:]...)
		spcIndex = reg.FindStringIndex(string(s2))
	}
	return string(s2)
}

var charStr = "zxcvbnmlkjhgfdsaqwertyuiopQWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func RandStr(length int) string {
	strArr := make([]string, 0)
	//长度为几就循环几次
	for i := 0; i < length; i++ {
		//产生0-61的数字
		number := rand.Intn(62)
		//将产生的数字通过length次承载到sb中
		strArr = append(strArr, string(charStr[number]))
	}
	//将承载的字符转换成字符串
	return strings.Join(strArr, "")
}

func DecodeUnicode(str string) string {
	sUnicodev := strings.Split(str, "\\u")
	var c string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		c += fmt.Sprintf("%c", temp)
	}
	return c
}
