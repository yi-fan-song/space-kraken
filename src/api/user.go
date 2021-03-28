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
	"net/http"
	"strings"
	"time"
)

// CreateUser is the response model for /users/:username/token
type CreateUser struct {
	Token string `json:"token"`
	User  struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Credits   int64     `json:"credits"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `json:"user"`
}

// CreateAccount creates an account with username
func (c Client) CreateAccount(username string) (token string, err error) {
	c.l.Printf("Creating an account with username %s...", username)

	url := "https://api.spacetraders.io/users/:username/token"
	url = strings.Replace(url, ":username", username, 1)

	var user CreateUser
	err = c.Do(url, http.MethodPost, nil, &user)
	if err != nil {
		c.l.Error("Creating account failed: ", err)
		return
	}
	return user.Token, err
}
