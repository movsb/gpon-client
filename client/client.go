package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// GponClient ...
type GponClient struct {
	ip     string
	cookie *http.Cookie
	token  string
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

	client := &GponClient{
		ip:     ip,
		cookie: sysauthCookie,
	}

	resp = client.mustGet(client.settingURL("status"))
	defer resp.Body.Close()
	source, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("cannot get token: %v", err)
	}
	reToken := regexp.MustCompile(`token: '([^']+)'`)
	matches := reToken.FindStringSubmatch(string(source))
	if len(matches) != 2 {
		log.Fatalf("cannot get token: %v", err)
	}
	client.token = matches[1]
	return client
}

func (c *GponClient) settingURL(name string) string {
	return fmt.Sprintf("http://%s/cgi-bin/luci/admin/settings/%s", c.ip, name)
}

func (c *GponClient) mustGet(u string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		log.Fatalf("cannot get url: %s: %v\n", u, err)
	}
	req.AddCookie(c.cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("cannot get url: %s: %v\n", u, err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("cannot get url: http status != 200: %v\n", resp.Status)
	}
	return resp
}

func (c *GponClient) mustPostForm(u string, data map[string]interface{}) *http.Response {
	values := url.Values{}
	for key, value := range data {
		values.Set(key, fmt.Sprint(value))
	}
	values.Set("token", c.token)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatalf("cannot post url: %s: %v\n", u, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(c.cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("cannot post url: %s: %v\n", u, err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("cannot post url: http status != 200: %v\n", resp.Status)
	}
	return resp
}

func (c *GponClient) mustGetJSON(out interface{}, url string) {
	resp := c.mustGet(url)
	defer resp.Body.Close()
	if out != nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(out); err != nil {
			log.Fatalf("cannot unmarshal json: %v\n", err)
		}
	}
}

func (c *GponClient) mustPostFormGetJSON(out interface{}, url string, data map[string]interface{}) {
	resp := c.mustPostForm(url, data)
	defer resp.Body.Close()
	if out != nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(out); err != nil {
			log.Fatalf("cannot unmarshal json: %v\n", err)
		}
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

// CreatePortMapping creates a port mapping.
//
// protocol: TCP, UDP, BOTH
func (c *GponClient) CreatePortMapping(name string, protocol string, outerPort int, innerIP string, innerPort int) {
	var ret RetVal
	c.mustPostFormGetJSON(&ret, c.settingURL("pmSetSingle"), map[string]interface{}{
		"op":       "add",
		"srvname":  name,
		"client":   innerIP,
		"protocol": protocol,
		"exPort":   outerPort,
		"inPort":   innerPort,
	})
	if ret.RetVal != 0 {
		log.Fatalf("cannot create port mapping: %v\n", ret.RetVal)
	}
}

// EnablePortMapping ...
func (c *GponClient) EnablePortMapping(name string, enable bool) {
	var ret RetVal
	op := "enable"
	if !enable {
		op = "disable"
	}
	c.mustPostFormGetJSON(&ret, c.settingURL("pmSetSingle"), map[string]interface{}{
		"op":      op,
		"srvname": name,
	})
	if ret.RetVal != 0 {
		log.Fatalf("cannot enable/disable port mapping: %v\n", ret.RetVal)
	}
}

// DeletePortMapping ...
func (c *GponClient) DeletePortMapping(name string) {
	var ret RetVal
	c.mustPostFormGetJSON(&ret, c.settingURL("pmSetSingle"), map[string]interface{}{
		"op":      "del",
		"srvname": name,
	})
	if ret.RetVal != 0 {
		log.Fatalf("cannot delete port mapping: %v\n", ret.RetVal)
	}
}

// GetGatewayInfo ...
func (c *GponClient) GetGatewayInfo() GatewayInfo {
	var gw _GatewayInfo
	c.mustGetJSON(&gw, c.settingURL(`gwinfo?get=part`))
	return gw.ToGatewayInfo()
}
