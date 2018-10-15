package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/devplayg/test_http_compress"
	"io/ioutil"
	"net/http"
	"os"
)

// Flag set
var fs *flag.FlagSet

func main() {
	// 옵션
	fs = flag.NewFlagSet("", flag.ExitOnError)
	var (
		compressed = fs.Bool("z", false, "body compressed")
		macCount = fs.Int("c", 500, "MAC Count")
	)
	fs.Usage = printHelp
	fs.Parse(os.Args[1:])

	// 가상 Mac 주소 생성
	macStr := test_http_compress.CreateFakeMac(*macCount)

	// 데이터 압축
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

func printHelp() {
	fmt.Println("[options]")
	fs.PrintDefaults()
}
