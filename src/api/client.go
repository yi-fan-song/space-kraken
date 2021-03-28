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

	"github.com/yi-fan-song/space-kraken/logger"
)

// Client is a client for requests to the space traders api
type Client struct {
	token    string
	username string

	httpClient *http.Client
	l          logger.Logger
}

// New creates a new api client
func New(token string, username string, httpClient *http.Client, l logger.Logger) Client {
	return Client{
		token:      token,
		username:   username,
		httpClient: httpClient,
		l:          l,
	}
}

// GetStatus gets the status of the game
func (c Client) GetStatus() (status GameStatus, err error) {
	c.l.Info("Fetching game Status...")

	url := "https://api.spacetraders.io/game/status"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.l.Error("Fetching failed: ", err)
		return
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.l.Error("Fetching failed: ", err)
		return
	}

	buf := make([]byte, 5000)
	n, err := res.Body.Read(buf)
	if err != nil && err != io.EOF {
		c.l.Error("Fetching failed: ", err)
		return
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			c.l.Error("Error occured but may not affect the output, ", err)
		}
	}()

	err = json.Unmarshal(buf[:n], &status)
	if err != nil {
		c.l.Error("Fetching failed: ", err)
		return
	}

	c.l.Info("Game status fetched:", status.Status)
	return
}

// Do "do"es a request
func (c Client) Do(url string, method string, body io.Reader, v interface{}) (err error) {
	c.l.Infof("Making a %s request to url: %s", method, url)

	var req *http.Request
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
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
			c.l.Error("Error occured but may not affect the output, ", err)
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
