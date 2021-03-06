// +build terminal

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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/komkom/toml"
	"github.com/yi-fan-song/space-kraken/api"
	"github.com/yi-fan-song/space-kraken/logger"
)

// Settings is struct for parsing settings
type Settings struct {
	User struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	} `json:"user"`
}

var l logger.Logger
var httpClient http.Client
var gameClient api.Client

var settings Settings

func init() {
	l = logger.New(ioutil.Discard)
	// l = logger.New(os.Stdout)

	if err := loadSettings("settings.toml", &settings); err != nil {
		l.Error(err)
		os.Exit(1)
	}

	httpClient = http.Client{Timeout: time.Minute}
	gameClient = api.New(settings.User.Token, settings.User.Name, &httpClient, l)
}

func loadSettings(filename string, settings *Settings) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(toml.New(bytes.NewBuffer(f)))
	if err := dec.Decode(&settings); err != nil {
		return err
	}

	return nil
}

func waitWhileOffline() {
	fmt.Println("Checking api status")
	for {
		status, err := gameClient.GetStatus()
		if err != nil {
			l.Error(err)
		}
		if status.Status == api.OkStatusMessage {
			fmt.Println("api is online, the game is available to play")
			break
		}

		fmt.Println("Waiting for api to come online")
		time.Sleep(time.Second * 30)
	}
}

func main() {
	waitWhileOffline()
	startPrompts()
}
