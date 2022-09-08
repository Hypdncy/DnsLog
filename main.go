package main

import (
	"DnsLog/Core"
	"DnsLog/Dns"
	"DnsLog/Http"
	"encoding/json"
	"io/ioutil"
	"log"
)

//GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-w -s" main.go

func main() {

	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	err = json.Unmarshal(content, &Core.Config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	go Dns.ListingDnsServer()
	Http.ListingHttpManagementServer()

}
