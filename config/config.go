package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"strings"
	"time"
)

type Config struct {
	Token    string   `json:"token,omitempty"`
	Domains  *Domains `json:"domains,omitempty"`
	IPSource string   `json:"ipSource,omitempty"`
	TTL      int      `json:"ttl,omitempty"`
	Cache    bool     `json:"cache,omitempty"`
}

type Domains struct {
	IPv4 []string `json:"ipv4,omitempty"`
	IPv6 []string `json:"ipv6,omitempty"`
}

func Load() (*Config, *time.Time, error) {
	config := &Config{
		Domains: new(Domains),
		Cache:   true,
	}

	configPath := flag.String("config", "", "config file path")
	flag.StringVar(&config.Token, "token", "", "API Token or Secret Key")
	ipv4Domains := flag.String("ipv4", "", "IPv4 domain list")
	ipv6Domains := flag.String("ipv6", "", "IPv6 domain list")
	flag.StringVar(&config.IPSource, "ip-source", "", "the source of IP address")
	flag.IntVar(&config.TTL, "ttl", 0, "TTL for DNS record")
	flag.Parse()

	if *configPath == "" {
		if config.Token == "" {
			return nil, nil, errors.New("Config error: Config file or API token are required.")
		}

		if *ipv4Domains != "" {
			config.Domains.IPv4 = strings.Split(*ipv4Domains, " ")
		}
		if *ipv6Domains != "" {
			config.Domains.IPv6 = strings.Split(*ipv6Domains, " ")
		}
		return config, nil, nil
	}

	file, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, nil, err
	}
	info, err := os.Stat(*configPath)
	if err != nil {
		return nil, nil, err
	}
	modTime := info.ModTime()

	return config, &modTime, json.Unmarshal(file, config)
}
