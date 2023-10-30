package main

import (
	"fmt"
	"log"

	"github.com/bionicosmos/ddns/address"
	"github.com/bionicosmos/ddns/cache"
	"github.com/bionicosmos/ddns/config"
	"github.com/bionicosmos/ddns/dns"
)

func main() {
	config, modTime, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	enableIPv4 := len(config.Domains.IPv4) != 0
	enableIPv6 := len(config.Domains.IPv6) != 0

	if config.Cache {
		cache, err := cache.Load()
		if err != nil {
			log.Fatal(err)
		}

		updated, err := cache.Update(config, modTime, enableIPv4, enableIPv6)
		if err != nil {
			log.Fatal(err)
		}
		if !updated {
			fmt.Println("Nothing to do :)")
			return
		}
	}

	if enableIPv4 {
		ipv4Address, err := address.Get(config, address.IPv4)
		if err != nil {
			log.Fatal(err)
		}
		for _, domain := range config.Domains.IPv4 {
			log.Printf("Updating %v to %v...", domain, ipv4Address)
			if err := dns.Update("A", domain, ipv4Address, config.Token, config.TTL); err != nil {
				log.Print(err)
			}
		}
	}

	if enableIPv6 {
		ipv6Address, err := address.Get(config, address.IPv6)
		if err != nil {
			log.Fatal(err)
		}
		for _, domain := range config.Domains.IPv6 {
			log.Printf("Updating %v to %v...", domain, ipv6Address)
			if err := dns.Update("AAAA", domain, ipv6Address, config.Token, config.TTL); err != nil {
				log.Print(err)
			}
		}
	}
}
