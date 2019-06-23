package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type GponClient struct {
	ip string
	// The sysauth login cookie
	cookie *http.Cookie
}

// MustDial tries to log in and then returns a GponClient
// that can be used to modify configurations.
func MustDial(ip string, username string, password string) *GponClient {
	httpClient := &http.Client{
		// disable redirects to get cookies
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	data := url.Values{}
	data.Set("username", username)
	data.Set("psd", password)
	resp, err := httpClient.PostForm(fmt.Sprintf("http://%s/cgi-bin/luci", ip), data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var sysauthCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "sysauth" {
			sysauthCookie = cookie
			break
		}
	}
	if sysauthCookie == nil {
		log.Fatal("login failed.")
	}

	return &GponClient{
		ip:     ip,
		cookie: sysauthCookie,
	}
}

func (c *GponClient) settingURL(name string) string {
	return fmt.Sprintf("http://%s/cgi-bin/luci/admin/settings/%s", c.ip, name)
}

func (c *GponClient) mustGet(url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("cannot get url: %s: %v\n", url, err)
	}
	req.AddCookie(c.cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("cannot get url: %s: %v\n", url, err)
	}
	return resp
}

func (c *GponClient) mustGetJSON(out interface{}, url string) {
	resp := c.mustGet(url)
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(out); err != nil {
		log.Fatalf("cannot unmarshal json: %v\n", err)
	}
}

// ListPortMappings list port mappings.
func (c *GponClient) ListPortMappings() []*PortMappingRule {
	rawObject := map[string]interface{}{}
	c.mustGetJSON(&rawObject, c.settingURL("pmDisplay"))
	outRules := make([]*PortMappingRule, 0, len(rawObject))
	for key, value := range rawObject {
		if !strings.HasPrefix(key, "pmRule") {
			continue
		}
		bys, err := json.Marshal(value)
		if err != nil {
			log.Fatalf("cannot unmarshal json: %v\n", err)
		}
		var rule PortMappingRule
		if err := json.Unmarshal(bys, &rule); err != nil {
			log.Fatalf("cannot unmarshal json: %v\n", err)
		}
		outRules = append(outRules, &rule)
	}
	return outRules
}
