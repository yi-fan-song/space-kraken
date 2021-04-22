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

// Package database provides a client for sqlite database operations.
package database

import (
	"github.com/yi-fan-song/space-kraken/log"
	"gorm.io/gorm"
)

// Client is the client for database operations.
type Client struct {
	db     *gorm.DB
	logger log.Logger
}

// New creates a new database client.
func New(db *gorm.DB, logger log.Logger) Client {
	return Client{
		db:     db,
		logger: logger,
	}
}

// MigrateModels creates all the necessary tables in the database
func (c Client) MigrateModels() {
	c.db.AutoMigrate(&User{})
}
