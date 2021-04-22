/**
 * Copyright (C) 2021 Yi Fan Song <yfsong00@gmail.com>
 *
 * This file is part of space-kraken.
 *
 * space-kraken is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * space-kraken is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with space-kraken.  If not, see <https://www.gnu.org/licenses/>.
 **/

package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/yi-fan-song/space-kraken/log"
)

const (
	// OkStatusMessage represents the message that the api will return when it is online
	OkStatusMessage = "spacetraders is currently online and available to play"

	BaseUrl = "https://api.spacetraders.io"
)

// Client is a client for the space traders api
type Client struct {
	username string
	token    string

	httpClient *http.Client
	logger     log.Logger
}

// New creates a new api client
func New(username string, token string, httpClient *http.Client, logger log.Logger) Client {
	return Client{
		username:   username,
		token:      token,
		httpClient: httpClient,
		logger:     logger,
	}
}

// Headers is an alias for a map string->string
type Headers map[string]string

// Do does a request and parses the response into v, type of v should
// correspond to expected response.
//
// If the api returns an error, a corresponding error will be returned.
func (c Client) Do(url string, method string, body io.Reader, headers Headers, v interface{}) (err error) {
	c.logger.Infof("Making a %s request to url: %s", method, url)

	var req *http.Request
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	var res *http.Response
	res, err = c.httpClient.Do(req)
	if err != nil {
		return
	}

	var buf []byte
	buf, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			c.logger.Error("Error occured but may not affect the output, ", err)
		}
	}()

	var apiErr ResponseError
	err = json.Unmarshal(buf, &apiErr)
	if err != nil {
		return
	}

	if apiErr.Error.Code != 0 {
		err = apiErr.ToClientError()
		return
	}

	err = json.Unmarshal(buf, v)
	return
}

func createAuthHeader(token string) Headers {
	return Headers{"Authorization": "Bearer " + token}
}

// HasAuth returns true if username and token is set
func (c *Client) HasAuth() bool {
	return c.username != "" && c.token != ""
}

// SetAuth sets username and token
func (c *Client) SetAuth(username string, token string) {
	c.username = username
	c.token = token
}
