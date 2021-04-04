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

package database

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Token    string
}

func (c Client) FetchUser() User {
	var user User
	tx := c.db.First(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return User{}
	}
	if tx.Error != nil {
		c.logger.Error(tx.Error)
	}
	return user
}

func (c Client) UpdateOrCreateUser(username string, token string) error {
	var user User
	tx := c.db.First(&user)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		c.db.Create(&User{
			Username: username,
			Token:    token,
		})
	} else if tx.Error != nil {
		c.logger.Error(tx.Error)
		return tx.Error
	} else {
		c.db.Model(&user).Updates(User{
			Username: username,
			Token:    token,
		})
	}
	return nil
}
