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
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/komkom/toml"
	"github.com/yi-fan-song/space-kraken/api"
	"github.com/yi-fan-song/space-kraken/database"
	"github.com/yi-fan-song/space-kraken/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	configDir  = "/etc/space-kraken"
	configPath = configDir + "/settings.toml"
	dataPath   = configDir + "/data"
	logPath    = configDir + "/latest.log"
)

var (
	logger log.Logger

	httpClient http.Client
	gameClient api.Client
	dbClient   database.Client

	settings Settings
)

func init() {
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(configDir, 0755)
	}

	err := loadSettings(configPath, &settings)
	if err != nil {
		panic(err)
	}

	logfile := createLogfile(logPath)
	logger = log.New(logfile, logfile, log.Config{
		ErrorColor: log.Red,
		InfoColor:  log.Cyan,
		UseColor:   settings.Logging.Color,
	})

	gLogger := gormLogger.New(PrintFormater{logger}, gormLogger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      gormLogger.Warn,
		Colorful:      settings.Logging.Color,
	})

	db, err := gorm.Open(sqlite.Open(dataPath), &gorm.Config{Logger: gLogger})
	if err != nil {
		panic("failed to connect database")
	}

	dbClient = database.New(db, logger)
	dbClient.MigrateModels()
	user := dbClient.FetchUser()

	httpClient = http.Client{Timeout: time.Minute}
	gameClient = api.New(user.Username, user.Token, &httpClient, logger)
}

func main() {
	waitWhileOffline()
	startPrompts()
}

// Settings is struct for parsing settings
type Settings struct {
	Logging struct {
		Color bool   `json:"color"`
		Level string `json:"level"`
	} `json:"logging"`
}

type PrintFormater struct {
	log.Logger
}

func (l PrintFormater) Printf(format string, a ...interface{}) {
	l.Infof(format, a...)
}

func createLogfile(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("Could not create log file: %s.\n", err)
		return nil
	}

	return f
}

func loadSettings(path string, settings *Settings) error {
	f, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(path)
			if err != nil {
				return err
			}
			if err := f.Chmod(0644); err != nil {
				return err
			}
			return loadSettings(f.Name(), settings)
		}
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
		status, err := gameClient.FetchStatus()
		if err != nil {
			logger.Error(err)
		}
		if status.Status == api.OkStatusMessage {
			fmt.Println("api is online, the game is available to play")
			break
		}

		fmt.Println("Waiting for api to come online")
		time.Sleep(time.Second * 30)
	}
}
