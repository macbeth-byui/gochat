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
* go run server/server.go

To start client:
* go run client/client.go

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

Rust Dependencies:
* wasm-bindgen = "0.2.100"
* js-sys = "0.3"
* serde = { version = "1.0", features = ["derive"] }
* serde_json = "1.0"
* regex = "1.11.1"
* wasm-bindgen-futures = "0.4.50"
* serde-wasm-bindgen = "0.6.5"

JavaScript Dependencies:
* "react": "^19.1.0",
* "react-dom": "^19.1.0",
* "react-scripts": "5.0.1",
* "react-virtualized": "^9.22.6",
* "wasm-gospel-search": "file:../pkg",

## Useful Websites to Learn More

I found these websites useful in developing this software:

* [React Virtualized](https://github.com/bvaughn/react-virtualized)
* [Rust and WebAssembly](https://rustwasm.github.io/docs/book/)
* [React](https://react.dev/reference/react)
* [Serde JSON](https://github.com/serde-rs/json)

## Future Work

The following items I plan to fix, improve, and/or add to this project in the future:

* [x] Refactor the code for readability, efficiency, and documentation
* [x] Host on a website
