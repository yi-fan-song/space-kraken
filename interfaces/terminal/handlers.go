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
	"fmt"
	"os"
)

func handleCmd(args []string) {
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "status":
		handleStatus()
	case "account":
		handleAccount(args[1:])
	case "exit":
		os.Exit(0)
	}
}

func handleStatus() {
	status, err := gameClient.FetchStatus()
	if err != nil {
		fmt.Println("Failed to fetch status:", err)
	} else {
		fmt.Println(status.Status)
	}
}

func handleAccount(args []string) {
	switch args[0] {
	case "create":
		if len(args[1:]) < 1 {
			fmt.Println("There are not enough arguments")
			break
		}
		username := args[1]

		user := dbClient.FetchUser()
		if user.Username != "" && user.Token != "" {
			fmt.Printf("You are already logged in as %s, are you sure you want to create another account?\n", user.Username)
			fmt.Printf("The token will be lost, cancel this command and use \"account token\" to retrive your token.\n")
			if !promptForYes("Confirm [yes/no]?", nil) {
				return
			}
		}

		token, err := gameClient.CreateAccount(username)
		if err != nil {
			fmt.Println("Could not create that account: ", err)
			return
		}

		err = dbClient.UpdateOrCreateUser(username, token)
		if err != nil {
			fmt.Printf("The account has been created but failed to save: %s.\n", err)
			fmt.Printf("Save this token to try again with \"account login\".\n")
			fmt.Printf("Token: %s", token)
			return
		}

		gameClient.SetAuth(username, token)
		fmt.Printf("Created account! Username: %s Token: %s\n", username, token)
		fmt.Println("Make sure to keep that token safe")

	case "login":
		if len(args[1:]) < 2 {
			fmt.Println("There are not enough arguments")
			break
		}
		username := args[1]
		token := args[2]
		gameClient.SetAuth(username, token)

		if _, err := gameClient.FetchAccount(); err != nil {
			fmt.Printf("Could not verify username/token pair: %s\n", err.Error())
			break
		}

		user := dbClient.FetchUser()
		if user.Username != "" && user.Token != "" {
			fmt.Printf("You are already logged in as %s, are you sure you want to login as %s?\n", user.Username, username)
			fmt.Printf("The token will be lost, cancel this command and use \"account token\" to retrive your token.\n")
			if !promptForYes("Confirm [yes/no]?", nil) {
				gameClient.SetAuth(user.Username, user.Token)
				return
			}
		}

		if err := dbClient.UpdateOrCreateUser(username, token); err != nil {
			fmt.Printf("Failed to save the account: %s.\n", err)
			return
		}
		fmt.Printf("Successfully logged in as %s.\n", username)

	case "token":
		user := dbClient.FetchUser()
		if user.Token == "" {
			fmt.Printf("You haven't logged in yet.\n")
			return
		}
		fmt.Printf("Logged in with username: %s, token: %s.\n", user.Username, user.Token)
	}
}
