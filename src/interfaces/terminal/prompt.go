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
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
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
	var user User
	result := db.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		fmt.Println("You haven't set an username and/or token, you can create an account by typing \"account create <username>\"")
		fmt.Println("If you have an account, you can log in with \"account login <username> <token>\"")
	} else {
		fmt.Printf("You've logged in as %s\n", user.Username)
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

			token, err := gameClient.CreateAccount(cmd[2])
			if err != nil {
				fmt.Println("Could not create that account: ", err)
				break
			}

			var user User
			r := db.First(&user)
			if errors.Is(r.Error, gorm.ErrRecordNotFound) {
				db.Create(&User{
					Username: cmd[2],
					Token:    token,
				})
				fmt.Printf("The account was saved, you are now logged in as %s.\n", cmd[2])
			} else {
				ret := promptForYes(
					fmt.Sprintf("You already have account with username %s, do you want to overwrite it?", user.Username),
					nil,
				)
				if ret {
					db.Model(&user).Updates(User{
						Username: cmd[2],
						Token:    token,
					})
					fmt.Printf("The account was saved, you are now logged in as %s.\n", cmd[2])
				} else {
					fmt.Println("The account was created, but not saved, you should save the token.")
				}
			}

			fmt.Printf("Created account! Username: %s Token: %s\n", cmd[2], token)
			fmt.Println("Make sure to keep that token safe")

		case "login":
			if len(cmd[2:]) < 2 {
				fmt.Println("There are not enough arguments")
				break
			}
		}
	case "exit":
		os.Exit(0)
	}
}

func promptForYes(message string, callback func()) bool {
	fmt.Print(message)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	args := strings.Split(text, " ")
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
	fmt.Print(message)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	return strings.Split(text, " ")
}
