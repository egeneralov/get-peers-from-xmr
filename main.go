package main


import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)


type rpcAnswer struct {
	GrayList []GrayList `json:"gray_list"`
}
type GrayList struct {
	Host     string `json:"host"`
	ID       uint64 `json:"id"`
	IP       int64  `json:"ip"`
	LastSeen uint64 `json:"last_seen"`
	Port     int64  `json:"port"`
}

type resultType struct {
	Host string
	Port string
}

func (self *resultType) String() string {
	return fmt.Sprintf("%s:%s", self.Host, self.Port)
}

func main() {
	var LocalNodeRpcUrl string
	flag.StringVar(&LocalNodeRpcUrl, "rpcurl", "http://127.0.0.1:18081", "url to daemon rpc")
	var verifyRequired bool
	flag.BoolVar(&verifyRequired, "verify", false, "check if port are available")
	flag.Parse()

	res, err := get(LocalNodeRpcUrl)
	if err != nil {
		panic(err)
	}

	for _, hp := range res {
		if verifyRequired {
			if raw_connect(hp.Host, hp.Port) {
				fmt.Println("alive:", hp.String())
			} else {
				fmt.Println("fail:", hp.String())
			}
		} else {
			fmt.Println(hp.String())
		}
	}
}

func raw_connect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		defer conn.Close()
		return true
	}
	return false
}

func InttoIP4(ipInt int64) string {
	// https://www.socketloop.com/tutorials/golang-convert-decimal-number-integer-to-ipv4-address
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func get(url string) ([]resultType, error) {
	var result []resultType
	var v rpcAnswer

	resp, err := http.Get(url + "/get_peer_list")
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return result, err
	}

	for _, a := range v.GrayList {
		r := resultType{
			Host: InttoIP4(a.IP),
			Port: fmt.Sprintf("%d", a.Port),
		}
		result = append(result, r)
	}

	return result, nil
}
