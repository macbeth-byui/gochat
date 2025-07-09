## Overview

**Project Title**

GoChat

**Project Description**

TCP Chat Server using Go

**Project Goals**

Learn the Go Language

## Instructions for Build and Use

This development version of GoChat is hardcoded to use 127.0.0.1:5001

To start server:
* `go run server/server.go`

To start client:
* `go run client/client.go`

Commands for use by the client:
* `USER|name` - Set the user name in the server directory
* `LIST` - Get a list of all active users in the server
* `CHAT_REQ|name|msg` - Send a chat message to a user (client must have a user name set)
* `BCAST_REQ|msg` - Send a change message to all users (client must have a user name set)
* `exit` - Disconnect from the server and exit the client

Potential Responses:
* `CHAT_RSP|name|msg` - Chat message from the specified name
* `OK|notes` - Previous command was successful
* `ERROR|notes` - Previous command was unsuccessful

To stop the server, press `CTRL-C`.

## Development Environment 

To recreate the development environment, you need the following software and/or libraries with the specified versions:

* Go 1.24.4

## Useful Websites to Learn More

I found these websites useful in developing this software:

* [Go By Example](https://gobyexample.com/)
* [Socket Programming in Go](https://www.kelche.co/blog/go/socket-programming/)
* [Install Go](https://go.dev/doc/install)

## Future Work

The following items I plan to fix, improve, and/or add to this project in the future:

* [ ] Add a GUI front-end

