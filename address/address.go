package address

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/bionicosmos/ddns/config"
)

type Version uint

const (
	IPv4 Version = iota
	IPv6
)

func Get(config *config.Config, version Version) (address string, err error) {
	if config.IPSource == "lan" {
		address, err = getLAN(version)
	} else {
		address, err = getWAN(version)
	}
	if err != nil {
		err = errors.New("Address error: Fail to get addresses.")
	}
	return
}

func getLAN(version Version) (string, error) {
	address := ""
	if version == IPv4 {
		address = "1.1.1.1:53"
	} else {
		address = "[2606:4700:4700::1111]:53"
	}

	conn, err := net.Dial("udp", address)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}

func getWAN(version Version) (string, error) {
	url := ""
	if version == IPv4 {
		url = "http://api-ipv4.ip.sb/ip"
	} else {
		url = "http://api-ipv6.ip.sb/ip"
	}

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	return strings.TrimSuffix(string(body), "\n"), err
}
