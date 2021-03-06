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
	"os"
	"strings"
	"time"

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

	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Minute,
		}
	}

	if l == nil {
		l = logger.New(os.Stdout)
	}

	return Client{
		token:      token,
		username:   username,
		httpClient: httpClient,
		l:          l,
	}
}

// GetStatus gets the status of the game
func (c Client) GetStatus() (status GameStatus, err error) {
	c.l.Log("Fetching game Status...")

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

	c.l.Log("Game status fetched:", status.Status)
	return
}

// CreateAccount creates an account with username
func (c Client) CreateAccount(username string) (token string, err error) {
	c.l.Logf("Creating an account with username %s...", username)

	url := "https://api.spacetraders.io/users/:username/token"
	url = strings.Replace(url, ":username", username, 1)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		c.l.Error("Creating account failed: ", err)
		return
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.l.Error("Creating account failed: ", err)
		return
	}

	buf := make([]byte, 5000)
	n, err := res.Body.Read(buf)
	if err != nil && err != io.EOF {
		c.l.Error("Creating account failed: ", err)
		return
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			c.l.Error("Error occured but may not affect the output, ", err)
		}
	}()

	user := CreateUser{}
	err = json.Unmarshal(buf[:n], &user)
	if err != nil {
		c.l.Error("Fetching failed: ", err)
		return
	}

	token = user.Token
	c.l.Logf("Created account with username %s...", username)
	return
}
