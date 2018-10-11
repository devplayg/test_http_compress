package main

import (
	"fmt"
	"github.com/icrowley/fake"
	"strings"
	"net/http"
	"io/ioutil"
	"bytes"
	"flag"
	"os"
	"compress/gzip"
)

// Flag set
var fs *flag.FlagSet

func main() {
	// 옵션
	fs = flag.NewFlagSet("", flag.ExitOnError)
	var (
		compressed         = fs.Bool("z", false, "body compressed")
	)
	fs.Usage = printHelp
	fs.Parse(os.Args[1:])

	// 가상 Mac 주소 생성
	macStr := createFakeMac(500)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(macStr)); err != nil {
		fmt.Errorf(err.Error())
		return
	}
	if err := gz.Close(); err != nil {
		fmt.Errorf(err.Error())
		return
	}
	//req, err := http.NewRequest("POST", q.host, &buf)

	// Request 객체 생성

	var req *http.Request
	var err error

	// 압축 옵션
	if *compressed {
		req, err = http.NewRequest("POST", "http://127.0.0.1", &buf)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Encoding", "gzip")
	} else {
		data := bytes.NewBuffer([]byte(macStr))
		req, err = http.NewRequest("POST", "http://127.0.0.1", data)
		if err != nil {
			panic(err)
		}
	}

	// 요청
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
	}

}

func createFakeMac(count int) string {
	var macList []string
	for i:=0; i<count; i++ {
		mac := fmt.Sprintf("AA:BB:CC:%s:%s:%s",
			fake.CharactersN(2),
			fake.CharactersN(2),
			fake.CharactersN(2),
		)
		macList = append(macList, mac)
	}
	return strings.ToUpper(strings.Join(macList, ","))
}

func printHelp() {
	fmt.Println("ekanited [options]")
	fs.PrintDefaults()
}