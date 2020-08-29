package server

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func userAgent() string {
	return "cddb" + "/" + "0.1" + " (https://github.com/defsub/cddb)"
}

func doGet(urlStr string) (*http.Response, error) {
	url, _ := url.Parse(urlStr)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func getJson(url string, result interface{}) error {
	resp, err := doGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(result); err != nil {
		return err
	}
	return nil
}
