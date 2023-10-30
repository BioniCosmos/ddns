package cache

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/bionicosmos/ddns/address"
	"github.com/bionicosmos/ddns/config"
)

type Cache struct {
	IPv4Address string
	IPv6Address string
	ModTime     *time.Time
}

const cachePath = "/tmp/ddns-cache.json"

func Load() (*Cache, error) {
	cacheFile, err := os.ReadFile(cachePath)
	if errors.Is(err, os.ErrNotExist) {
		newCacheFile, err := json.Marshal(Cache{})
		if err != nil {
			return nil, err
		}
		if err = os.WriteFile(cachePath, newCacheFile, 0644); err != nil {
			return nil, err
		}
		cacheFile = newCacheFile
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	cache := new(Cache)
	return cache, json.Unmarshal(cacheFile, cache)
}

func (cache *Cache) NeedUpdate(ipv4Address string, ipv6Address string, modTime *time.Time) bool {
	return !(cache.IPv4Address == ipv4Address && cache.IPv6Address == ipv6Address && cache.ModTime.Equal(*modTime))
}

func (cache *Cache) Update(config *config.Config, modTime *time.Time, enableIPv4 bool, enableIPv6 bool) (updated bool, err error) {
	ipv4Address := ""
	if enableIPv4 {
		ipv4Address, err = address.Get(config, address.IPv4)
		if err != nil {
			return
		}
	}

	ipv6Address := ""
	if enableIPv6 {
		ipv6Address, err = address.Get(config, address.IPv6)
		if err != nil {
			return
		}
	}

	if !cache.NeedUpdate(ipv4Address, ipv6Address, modTime) {
		return
	}

	updated = true
	cache = &Cache{
		IPv4Address: ipv4Address, IPv6Address: ipv6Address, ModTime: modTime,
	}
	cacheFile, err := json.Marshal(cache)
	if err != nil {
		return
	}
	err = os.WriteFile(cachePath, cacheFile, 0644)
	return
}
