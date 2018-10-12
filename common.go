package test_http_compress

import (
	"fmt"
	"github.com/icrowley/fake"
	"strings"
)

const TOPIC = "unisem/iot/mqtt"

func CreateFakeMac(count int) string {
	var macList []string
	for i:=1; i<=count; i++ {
		mac := fmt.Sprintf("FAKEMACFAKEMAC%d",
			100+i,
		)
		macList = append(macList, mac)
	}
	return strings.ToUpper(strings.Join(macList, ","))
}
func CreateFakeMac_old(count int) string {
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
