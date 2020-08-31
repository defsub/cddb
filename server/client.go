// Copyright (C) 2020 The cdmbz Authors.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func userAgent() string {
	return "cdmbz" + "/" + "0.1" + " (https://github.com/defsub/cdmbz)"
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
