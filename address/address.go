package address

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/bionicosmos/ddns/config"
	"golang.org/x/net/proxy"
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
		address, err = getWAN(version, config.Proxy)
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

func getWAN(version Version, p string) (string, error) {
	u := ""
	if version == IPv4 {
		u = "http://api-ipv4.ip.sb/ip"
	} else {
		u = "http://api-ipv6.ip.sb/ip"
	}

	client := http.DefaultClient
	if p != "" {
		proxyURL, err := url.Parse(p)
		if err != nil {
			return "", fmt.Errorf("Error parsing proxy URL: %w", err)
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return "", fmt.Errorf("Error creating proxy dialer: %w", err)
		}

		client = &http.Client{
			Transport: &http.Transport{
				Dial: dialer.Dial,
			},
		}
	}

	response, err := client.Get(u)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	return strings.TrimSuffix(string(body), "\n"), err
}
