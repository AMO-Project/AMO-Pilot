package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func GetPublicIP() (ip []byte) {
	rsp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}
	defer rsp.Body.Close()

	buf, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("[-] Error occured: %s\n", err)
		return nil
	}

	return net.ParseIP(string(bytes.TrimSpace(buf))).To4()
}
