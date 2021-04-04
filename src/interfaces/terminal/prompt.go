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
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startPrompts() {
	checkAccount()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		text := scanner.Text()

		cmd := strings.Split(text, " ")
		handleCmd(cmd)
	}
}

func checkAccount() {
	user := dbClient.FetchUser()
	if user.Username == "" && user.Token == "" {
		fmt.Println("You haven't set an username and/or token, you can create an account by typing \"account create <username>\"")
		fmt.Println("If you have an account, you can log in with \"account login <username> <token>\"")
	} else {
		fmt.Printf("You've logged in as %s\n", user.Username)
	}
}

func promptForYes(message string, callback func()) bool {
	args := promptAndWait(message)
	answer := strings.ToLower(args[0])

	if answer == "yes" || answer == "y" {
		if callback != nil {
			callback()
		}
		return true
	} else {
		return false
	}
}

func promptAndWait(message string) []string {
	fmt.Println(message)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	return strings.Split(text, " ")
}
