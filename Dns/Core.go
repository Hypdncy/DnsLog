package Dns

import (
	"DnsLog/Core"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"log"
	"net"
	"sync"
	"time"
)

var GlobalData = make(map[string][]LogInfo)

var rw sync.RWMutex

type LogInfo struct {
	Domain string `json:"domain"`
	Ip     string `json:"ip"`
	Time   int64  `json:"time"`
}

var L LogInfo

// ListingDnsServer 监听dns端口
func ListingDnsServer() {
	//if runtime.GOOS != "windows" && os.Geteuid() != 0 {
	//	log.Fatal("Please run as root")
	//}
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(Core.Config.DNS.Ip),
		Port: Core.Config.DNS.Port},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	log.Println("DNS Listing Start...")
	for {
		buf := make([]byte, 512)
		_, addr, _ := conn.ReadFromUDP(buf)
		var msg dnsmessage.Message
		if err := msg.Unpack(buf); err != nil {
			fmt.Println(err)
			continue
		}
		go serverDNS(addr, conn, msg)
	}
}

func serverDNS(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	if len(msg.Questions) < 1 {
		return
	}
	question := msg.Questions[0]
	var (
		queryDomain  = question.Name.String()
		queryType    = question.Type
		queryName, _ = dnsmessage.NewName(queryDomain)
		resource     dnsmessage.Resource
	)

	user, ok := Core.GetUser(queryDomain)
	if ok {
		L.Set(user, LogInfo{
			Domain: queryDomain,
			Ip:     addr.IP.String(),
			Time:   time.Now().Unix(),
		})
	}

	switch queryType {
	case dnsmessage.TypeA:
		resource = NewAResource(queryName, [4]byte{127, 0, 0, 1})
	default:
		resource = NewAResource(queryName, [4]byte{127, 0, 0, 1})
	}
	// send response
	msg.Response = true
	msg.Answers = append(msg.Answers, resource)
	Response(addr, conn, msg)
}

// Response return
func Response(addr *net.UDPAddr, conn *net.UDPConn, msg dnsmessage.Message) {
	packed, err := msg.Pack()
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := conn.WriteToUDP(packed, addr); err != nil {
		fmt.Println(err)
	}
}

func NewAResource(query dnsmessage.Name, a [4]byte) dnsmessage.Resource {
	return dnsmessage.Resource{
		Header: dnsmessage.ResourceHeader{
			Name:  query,
			Class: dnsmessage.ClassINET,
			TTL:   0,
		},
		Body: &dnsmessage.AResource{
			A: a,
		},
	}
}

func (d *LogInfo) Set(user string, data LogInfo) {
	rw.Lock()
	if GlobalData[user] == nil {
		GlobalData[user] = []LogInfo{data}
	} else {
		GlobalData[user] = append(GlobalData[user], data)
	}
	rw.Unlock()
	return
}

func (d *LogInfo) Get(user string) ([]LogInfo, bool) {
	rw.RLock()
	var res, ok = GlobalData[user]
	rw.RUnlock()
	return res, ok
}

func (d *LogInfo) Clear(user string) {
	GlobalData[user] = []LogInfo{}
	GlobalData["other"] = []LogInfo{}
}
