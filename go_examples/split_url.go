package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Enter your url here:")

	url := bufio.NewReader(os.Stdin)
	line, _ := url.ReadString('\n')
	paramStr := strings.Split(line, "?")[1]
	params := strings.Split(paramStr, "&")

	print("\n")
	for _, param := range params {
		percentSpl := strings.Split(param, "%")

		var strPara string
		if len(percentSpl) > 1 {
			for i, j := range percentSpl {
				if i == 0 {
					strPara += j
				} else {
					bl, _ := hex.DecodeString(j[:2])
					strung := string(bl)
					strPara += strung
					strPara += j[2:]
				}
			}
		} else {
			strPara = param
		}
		fmt.Println(strPara)
	}
}
