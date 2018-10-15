package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/devplayg/test_http_compress"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var client mqtt.Client
var fs *flag.FlagSet

func main() {
	// 옵션
	fs = flag.NewFlagSet("", flag.ExitOnError)
	var (
		clientId = fs.String("cid", "test-client-id", "Client ID")
		macCount = fs.Int("c", 500, "MAC Count")

	)
	fs.Usage = printHelp
	fs.Parse(os.Args[1:])

	// Client 생성
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(*clientId)
	opts.HTTPHeaders.Set("woosanheader", "woosanvalue")

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// MAC 문자열 생성
	macStr := test_http_compress.CreateFakeMac(*macCount)

	// 압축
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

	if token := client.Publish(test_http_compress.TOPIC, 0, false, buf.Bytes()); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	WaitForSignals()
}

func WaitForSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalCh:
		log.Println("Signal received, shutting down...")
	}
}

func handler(res http.ResponseWriter, req *http.Request) {
	macStr := test_http_compress.CreateFakeMac(500)
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

	if token := client.Publish(test_http_compress.TOPIC, 0, false, buf.Bytes()); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Fprintf(res, "hello %s", "abc")
}

func printHelp() {
	fmt.Println("[options]")
	fs.PrintDefaults()
}
