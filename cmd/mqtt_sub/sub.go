package main

import (
	"fmt"
	"github.com/devplayg/test_http_compress"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/icrowley/fake"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var client mqtt.Client

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	//var wg sync.WaitGroup
	//wg.Add(1)
	if token := client.Subscribe(test_http_compress.TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Println(string(msg.Payload()))
		fmt.Println("======================================")
		//if string(msg.Payload()) != "mymessage" {
		//	log.Fatalf("want mymessage, got %s", msg.Payload())
		//}
		//	wg.Done()
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	//data := CreateFakeMac(500)
	//if token := client.Publish(TOPIC, 0, false, data); token.Wait() && token.Error() != nil {
	//	log.Fatal(token.Error())
	//}
	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))
	//
	WaitForSignals()
}

func CreateFakeMac(count int) string {
	var macList []string
	for i := 0; i < count; i++ {
		mac := fmt.Sprintf("AA:BB:CC:%s:%s:%s",
			fake.CharactersN(2),
			fake.CharactersN(2),
			fake.CharactersN(2),
		)
		macList = append(macList, mac)
	}
	return strings.ToUpper(strings.Join(macList, ","))
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
	data := CreateFakeMac(500)
	if token := client.Publish(test_http_compress.TOPIC, 0, false, data); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	fmt.Fprintf(res, "hello %s", "abc")
}
