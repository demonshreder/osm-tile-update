package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func main() {
	dat, err := ioutil.ReadFile("links")
	check(err)
	strData := strings.Split(string(dat), "\n")
	// fmt.Println(strData)
	// resp, err := http.Get(strData[1] + "/dirty")
	resp, err := http.Get("http://httpbin.org/user-agent")
	check(err)
	fmt.Println(resp)
	fmt.Println(strData[1] + "/dirty")
	// for i := range strData {
	//
	// }
}
