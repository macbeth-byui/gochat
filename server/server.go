package main

/* GoChat Server - Allows clients to connect and send
   chat messages.  Manages the user name list. */

import (
	"bufio"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
)

// Directory mapped in both direction and protected between threads
type Directory struct {
	users map[string]net.Conn
	conns map[net.Conn]string
	mut sync.Mutex
}

// Send a message to a client
func send(conn net.Conn, msg string) {
	fmt.Println("[", conn.RemoteAddr(), "] <-", msg)
	conn.Write([]byte(msg + "\n"))
}

// Thread to manage a client.  Reference to the directory
// provided so that we can lock the mutex before using.
func client(conn net.Conn, directory *Directory) {
	fmt.Println("[", conn.RemoteAddr(), "] Connected.")
	for {
		// Read from the client
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			// Client disconnected
			directory.mut.Lock()
			user, exist := directory.conns[conn]
			if exist {
				delete(directory.users, user)
				delete(directory.conns, conn)
			}
			directory.mut.Unlock()
			fmt.Println("[", conn.RemoteAddr(), "] Disconnected.")
			return
		}
		// Parse the request form the client
		data = strings.TrimSpace(data)
		parts := strings.Split(data, "|")
		fmt.Println("[", conn.RemoteAddr(), "] ->", data)
		directory.mut.Lock()
		switch parts[0] {
			// USER|name
			case "USER":
				if len(parts) != 2 {
					send(conn,"ERROR|Missing name in USER")
				} else {
					_, exist := directory.users[parts[1]]
					if exist {
						send(conn,"ERROR|Name already exists in USER")
					} else {
						_, exist := directory.conns[conn]
						if exist {
							send(conn,"ERROR|Name already set for USER")
						} else {
							directory.users[parts[1]] = conn
							directory.conns[conn] = parts[1]
							send(conn,"OK|Name set in USER")
						}
					}
				}
			// CHAT_REQ|name|msg 
			case "CHAT_REQ":
				if len(parts) != 3 {
					send(conn,"ERROR|Missing user or message in CHAT_REQ")
				} else {
					local_user, exist := directory.conns[conn]
					if exist {
						remote_conn, exist := directory.users[parts[1]]
						if exist {
							if remote_conn != conn {
								send(remote_conn,"CHAT_RSP|"+local_user+"|"+parts[2])
								send(conn,"OK|Message sent with CHAT_RSP")
							} else {
								send(conn,"ERROR|Cannot send CHAT_REQ to self")
							}
						} else {
							send(conn,"ERROR|Invalid user in CHAT_REQ")
						}
					} else {
						send(conn,"ERROR|Cannot send CHAT_REQ without registering user name")
					}
				}
			// BCAST_REQ|msg
			case "BCAST_REQ":
				if len(parts) != 2{
					send(conn, "ERROR|Missing message in BCAST_REQ")
				} else {
					local_user, exist := directory.conns[conn]
					if exist {
						for _, remote_conn := range directory.users {
							if remote_conn != conn {
								send(remote_conn,"CHAT_RSP|"+local_user+"|"+parts[1])
							}
						}
						send(conn,"OK|Broadcast Messages Sent with CHAT_RSP")
					} else {
						send(conn,"ERROR|Cannot send BCAST_REQ without registering user name")
					}
				}
			// LIST
			case "LIST":
				if len(parts) != 1 {
					send(conn,"ERROR|Additional parameters added to LIST")
				} else {
					users := make([]string, 0)
					for user := range directory.users {
						users = append(users, user)
					}
					sort.Strings(users)
					send(conn,"OK|"+strings.Join(users, ","))
				}
			default:
				send(conn,"ERROR|Invalid Command")
		}
		directory.mut.Unlock()
	}
}

// Manage the listening socket for the server
func server(addr string) {
	// Create TCP server
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error starting peer server: ", err)
		return
	}
	fmt.Println("[", listen.Addr(), "] Server Started")
	// Create empty directory
	directory := Directory {
		users: make(map[string]net.Conn),
		conns: make(map[net.Conn]string),
	}
	for {
		// Accept clients and start thread for each one
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting client: ",err)
			continue
		}
		go client(conn, &directory)
	}
}

// Entry function for the program
func main() {
	server_addr := "127.0.0.1:5001"
	fmt.Println("[", server_addr, "] Starting Server (Ctrl-C to Exit)")
	server(server_addr)
}