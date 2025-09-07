package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ProxyListResponse struct {
	Count     int `json:"count"`
	Next      any `json:"next"`
	Previous  any `json:"previous"`
	ProxyList []struct {
		ID                    string    `json:"id"`
		Username              string    `json:"username"`
		Password              string    `json:"password"`
		ProxyAddress          string    `json:"proxy_address"`
		Port                  int       `json:"port"`
		Valid                 bool      `json:"valid"`
		LastVerification      time.Time `json:"last_verification"`
		CountryCode           string    `json:"country_code"`
		CityName              string    `json:"city_name"`
		AsnName               string    `json:"asn_name"`
		AsnNumber             int       `json:"asn_number"`
		HighCountryConfidence bool      `json:"high_country_confidence"`
		CreatedAt             time.Time `json:"created_at"`
	} `json:"results"`
}

type ProxyManager struct {
	Authorization string
	Proxies       []string
}

func NewProxyManager(auth string) *ProxyManager {
	return &ProxyManager{Authorization: auth, Proxies: []string{}}
}

func (pm *ProxyManager) FetchProxies() error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://proxy.webshare.io/api/v2/proxy/list/?mode=direct&page=1&page_size=25", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", pm.Authorization)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var proxyResp ProxyListResponse
	if err := json.NewDecoder(resp.Body).Decode(&proxyResp); err != nil {
		return err
	}
	var proxies []string
	for _, p := range proxyResp.ProxyList {
		if p.Valid {
			proxy := "http://" + p.Username + ":" + p.Password + "@" + p.ProxyAddress + ":" + strconv.Itoa(p.Port)
			proxies = append(proxies, proxy)
		}
	}
	pm.Proxies = proxies
	return nil
}

func (pm *ProxyManager) FetchAndGetRandomProxy() (string, error) {
	pm.FetchProxies()

	return pm.GetRandomProxy()
}

func (pm *ProxyManager) GetRandomProxy() (string, error) {
	if len(pm.Proxies) == 0 {
		return "", errors.New("no proxies available")
	}
	rand.Seed(time.Now().UnixNano())
	return pm.Proxies[rand.Intn(len(pm.Proxies))], nil
}
