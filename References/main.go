package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"flag"
)

/*
Notes:
./foo is a package
go build - Makes binary out of file you run command in, looks for go file within folder
sudo shutdown -s -h => For shutting down a server
ssh ubuntu@ip-address
vi (file name) - opens vm
go build
go run (file name).go OR ./foo "parameter"
go fmt - formats go files by proper indentation

Computer Communication: Open port on your side, establish connection to port on other side
*/

/*
 * Fill in the missing parts of this code to complete the client-server
 * implementation. You will need to complete the following high level steps:
 *   1) Parse the command line arguments, so that running "pinpong client"
 *      executes the client code, while running "pingpong server" executes
 *			the server code.
 *	 2) Complete the client function so that it sends a message to the server
 *      once every second for 100 seconds. You can hardcode the address and port
 *			of the server.
 *	 3)	Complete the server function. All it needs to do is check for messages
 *	    from the client and respond with its own message. The server should
 *			stop listenting after it has received 100 messages.
 *   4) Add the ability to specify custom client and server messages from the
 *      command line.
 */

//Main Variables
var message string = "Hello"

//Function Call
func main() {
	clientOrServer := flag.String("whichFunc", "filler", "client or server")
        ipAddress := flag.String("ipAddress", "169.229.50.175", "input ip address")
        port := flag.Int("port", 6521, "input port number")
	flag.Parse()

	parameter := *clientOrServer
	if parameter == "client" {
		fmt.Println("Client Function Called")
		client(*ipAddress, *port)
	} else if parameter == "server" {
		fmt.Println("Server Function Called")
		server(*port)
	} else {
		fmt.Println("Please enter either 'client' or 'server' to call the respective functions")
	}
}

/*
Connects to the server
Returns a connection and a handler error

TODO: no compilation errors, now test
*/
func client(ipAddress string, port int) {
	//Needs network and address string as parameters
	//Address String Format: IP Address:Port, use same port in server code
	//Connect to Pi server, not go server
        iPAddPort := fmt.Sprintf("%s:%d", ipAddress, port)
	clientConn, err := net.Dial("tcp", iPAddPort)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Passed net.dial")

	//Print out connection
	fmt.Fprintf(clientConn, "GET / HTTP/1.0\r\n\r\n")
	/*bufio
	"buffered io"
	https://golang.org/pkg/bufio/
	Reader is a struct w/ attributes:
	buf          []byte
  rd           io.Reader // reader provided by the client
  r, w         int       // buf read and write positions
  err          error
	*/
	reader := bufio.NewReader(clientConn)
	sendMsg := "Hello\n"
	numIters := 100
	for i := 0; i < numIters; i++ {
		fmt.Printf("(%d) Sending: %s", i, sendMsg)
		if _, err := clientConn.Write([]byte(sendMsg)); err != nil {
			fmt.Println("Send failed:", err) //Prints error message
			os.Exit(1)
		}

		/*
		ReadString Functionality:
	  1. Reads until 1st occurence of delim in input.
		fyi delimiter: character sequence used to specify boundary between separate, independent regions
		2. Returns String w/ data up to delimiter
		3. Returns err != nil iff the returned data does not end in delim
		*/
		recvMsg, err := reader.ReadString('\n') //Replaced err w/ _ b/c not using err
		if err != nil {
			fmt.Println("Error Message", err)
			return
		}
		fmt.Printf("(%d) Received: %s", i, recvMsg)

		/*
		Time Documentation: https://golang.org/pkg/time/
		Time Golang - time.sleep()
		*/
		time.Sleep(time.Second)
	}
	clientConn.Close()
}

func server(port int) {
	//Listen function creates server
        portString := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", portString) //Network, IP Address:Port (Must be same as client's)
	if err != nil {
		fmt.Println("Listen failed:", err)
		os.Exit(1)
	}

	serverConn, err := ln.Accept()
	if err != nil {
		fmt.Println("Accept failed:", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(serverConn)

	numIters := 100
	sendMsg := "Greetings\n"
	for i := 0; i < numIters; i++ {
		recvMsgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Receive failed:", err)
			os.Exit(1)
		}
		fmt.Printf("(%d) Received: %s", i, string(recvMsgBytes))
		fmt.Printf("(%d) Sending: %s\n", i, sendMsg)
		if _, err := serverConn.Write([]byte(sendMsg)); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}
	}
	serverConn.Close()
}

/*
TODO: Questions:
1. For ReadString why is 'n' acceptable but "n" not
2. What is a delimiter and how can we use it?
3. SSH Command Explained
*/
