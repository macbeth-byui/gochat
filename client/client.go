package main

/* GoChat Client - Connects to the server (defaulted to localhost) and then 
   accepts commands:

   USER|name - Provide a user name for the client.  Must be done to send/receive
               chat messages.  Cannot be changed.
   CHAT_REQ|name|msg - Send a chat to the user name
   BCAST_REQ|msg - Sends a chat to all users
   exit - Disconnect from server and exit the client
*/

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

/* Read on a socket connection.  If the connection fails, the 
   reader_exit_chan will be closed.  This is used by the client
   to detect if the socket disconnects. */
func reader(conn net.Conn, reader_exit_chan chan struct{}) {
	data_chan := make(chan string)

	// Start thread for reading from socket
	go func() {
		for {
			data, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				// Close all channels to stop loops
				fmt.Println("[", conn.RemoteAddr(), "] Server Disconnected.")
				close(data_chan)
				close(reader_exit_chan)
				return
			}
			// Remove newlines and send for processing
			// Maybe I didn't need this in a separate thread, 
			// but I'm learning about Go and thought more threads
			// was a better idea.
			data = strings.TrimSpace(data)
			data_chan <- data
		}
	}()

	for data := range data_chan {
		// Parse all messages from the server
		parts := strings.Split(data, "|")
		switch parts[0] {
			case "OK" :
				if len(parts) == 1 {
					fmt.Println("\033[32m"+parts[0]+"\033[0m")
				} else if len(parts) >= 2 {
					fmt.Println("\033[32m"+parts[0]+": "+parts[1]+"\033[0m")
				}
			case "ERROR" :
				if len(parts) == 1 {
					fmt.Println("\033[31m"+parts[0]+"\033[0m")
				} else if len(parts) >= 2 {
					fmt.Println("\033[31m"+parts[0]+": "+parts[1]+"\033[0m")
				}
			case "CHAT_RSP" :
				if len(parts) == 3 {
					fmt.Println("\033[34m["+parts[1]+"\033[0m]: "+parts[2])
				}
		}
	}
}

func client(addr string) {
	// Connect via TCP to the server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to peer server: ", err)
		return
	}
	fmt.Println("[", conn.RemoteAddr(), "] Connected to Server")

	// Start thread to read socket
	reader_exit_chan := make(chan struct{})
	go reader(conn, reader_exit_chan)

	// Start thread to read keyboard 
	reader := bufio.NewReader(os.Stdin)
	input_chan := make(chan string)

	go func() {
		for {
			data, err := reader.ReadString('\n')
			data = strings.TrimSpace(data)
			if err != nil {
				if data == "exit" {
					// Close the connection with the server which
					// will trigger the reader to stop.
					close(input_chan)
					conn.Close()
					return
				}
				if data != "" {
					// Ignore blank messages
					input_chan <- data
				}
			} else {
				// If reader closes
				close(input_chan)
				conn.Close()
				return
			}
			
		}
	}()

	outer: for {
		select {
			case <-reader_exit_chan:
				// If this closes then the connection was closed
				// by the reader.
				break outer
			case data, ok := <-input_chan:
				if ok {
					conn.Write([]byte(data + "\n"))
				} else {
					break outer
				}
		}
	}
	// Allow the connection to server to close before exiting.
	<- reader_exit_chan
}

// Entry function for the program
func main() {
	server_addr := "127.0.0.1:5001"
	fmt.Println("[", server_addr, "] Connecting to Server")
	client(server_addr)	

	
}