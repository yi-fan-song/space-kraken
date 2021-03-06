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
	if settings.User.Name == "ERR_NO_NAME" || settings.User.Token == "ERR_NO_TOKEN" {
		fmt.Println("You haven't set an username and/or token, you can create an account by typing \"account create <username>\"")
		fmt.Println("If you have an account, you can log in with \"account login <username> <token>\"")
	}
}

func handleCmd(cmd []string) {
	if len(cmd) == 0 {
		return
	}

	switch cmd[0] {
	case "ping":
		fmt.Println("pong")
	case "status":
		status, err := gameClient.GetStatus()
		if err != nil {
			fmt.Println("Failed to fetch status:", err)
		} else {
			fmt.Println(status.Status)
		}
	case "account":
		switch cmd[1] {
		case "create":
			if len(cmd[2:]) < 1 {
				fmt.Println("There are not enough arguments")
				break
			}

		case "login":
			if len(cmd[2:]) < 2 {
				fmt.Println("There are not enough arguments")
				break
			}
		}
	}

}
