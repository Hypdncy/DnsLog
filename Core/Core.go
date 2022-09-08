package Core

import (
	"strings"
)

var Config = struct {
	HTTP struct {
		Ip             string `json:"ip"`
		Port           int    `json:"port"`
		ConsoleDisable bool   `json:"consoleDisable"`
	} `json:"HTTP"`

	USER map[string]string `json:"USER"`
	DNS  struct {
		Ip     string `json:"ip"`
		Port   int    `json:"port"`
		Domain string `json:"domain"`
	} `json:"DNS"`
}{}

func GetUser(domain string) (string, bool) {
	var subDomain = strings.Replace(domain, "."+Config.DNS.Domain+".", "", -1)
	var subs = strings.Split(subDomain, ".")
	var user = subs[len(subs)-1]
	_, ok := Config.USER[user]
	if ok {
		return user, true
	}
	return user, false
}
