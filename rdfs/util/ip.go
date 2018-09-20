package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func GetPublicIP() (ip [4]byte) {
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return [4]byte{}
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return [4]byte{}
	}

	ip4byte := [4]byte{}
	copy(ip4byte[:], net.ParseIP(string(bytes.TrimSpace(buf))).To4())

	return ip4byte
}

func EncodeIP(rawIP [4]byte) string {
	return net.IPv4(rawIP[0], rawIP[1], rawIP[2], rawIP[3]).String()
}
