package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func main() {
	str := "4"
	pat := "^[1-4]$"

	if ok, _ := regexp.Match(pat, []byte(str)); ok {
		fmt.Println("匹配成功")
	} else {
		fmt.Println("匹配失败")
	}

	trueValue, _ := strconv.ParseInt("34230412834", 10, 32)
	fmt.Println(trueValue)
}
