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
)

// GameStatus is the struct for the get status response
type GameStatus struct {
	Status string `json:"status"`
}

// FetchStatus gets the status of the api
func (c Client) FetchStatus() (status GameStatus, err error) {
	c.logger.Info("Fetching game Status...")

	url := "https://api.spacetraders.io/game/status"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("Fetching failed: ", err)
		return
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Fetching failed: ", err)
		return
	}

	buf := make([]byte, 5000)
	n, err := res.Body.Read(buf)
	if err != nil && err != io.EOF {
		c.logger.Error("Fetching failed: ", err)
		return
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			c.logger.Error("Error occured but may not affect the output, ", err)
		}
	}()

	err = json.Unmarshal(buf[:n], &status)
	if err != nil {
		c.logger.Error("Fetching failed: ", err)
		return
	}

	c.logger.Info("Game status fetched:", status.Status)
	return
}
