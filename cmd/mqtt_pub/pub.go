package main

import (
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
	)
	fs.Usage = printHelp
	fs.Parse(os.Args[1:])

	// Client 생성
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(*clientId)
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))

	data := test_http_compress.CreateFakeMac(500)
	if token := client.Publish(test_http_compress.TOPIC, 0, false, data); token.Wait() && token.Error() != nil {
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
	data := test_http_compress.CreateFakeMac(500)
	if token := client.Publish(test_http_compress.TOPIC, 0, false, data); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Fprintf(res, "hello %s", "abc")
}

func printHelp() {
	fmt.Println("[options]")
	fs.PrintDefaults()
}
