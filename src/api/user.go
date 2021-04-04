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
	"fmt"
	"net/http"
	"strings"
)

// InnerUser is the type embeded in user endpoint responses
type InnerUser struct {
	Username string                   `json:"username"`
	Credits  int64                    `json:"credits"`
	Ships    []map[string]interface{} `json:"ships"`
	Loans    []map[string]interface{} `json:"loans"`
}

// CreatedUser is the response model for /users/:username/token
type CreatedUser struct {
	Token string    `json:"token"`
	User  InnerUser `json:"user"`
}

// FetchedUser is the response model for /users/:username
type FetchedUser struct {
	User InnerUser `json:"user"`
}

// CreateAccount creates an account with username
func (c Client) CreateAccount(username string) (token string, err error) {
	c.logger.Infof("Creating an account with username %s...", username)

	url := BaseUrl + "/users/:username/token"
	url = strings.Replace(url, ":username", username, 1)

	var user CreatedUser
	err = c.Do(url, http.MethodPost, nil, nil, &user)
	if err != nil {
		c.logger.Error("Creating account failed: ", err)
		return
	}
	return user.Token, err
}

// FetchAccount fetches an account with the username and token
func (c Client) FetchAccount() (user FetchedUser, err error) {
	if !c.HasAuth() {
		err = fmt.Errorf("Client without auth only supports account creation and status")
		return
	}

	c.logger.Infof("Fetching the account with username %s...", c.username)

	url := BaseUrl + "/users/:username"
	url = strings.Replace(url, ":username", c.username, 1)

	headers := createAuthHeader(c.token)
	err = c.Do(url, http.MethodGet, nil, headers, &user)
	if err != nil {
		c.logger.Error("Fetching user failed: ", err)
	}

	return
}
