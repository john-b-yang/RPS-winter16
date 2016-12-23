/*
 * Designated import statements and necessary packages
 */
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
Server ipAddress and port can be parameters
Client port will be assigned.
For two people to play, we want to call the server function once and client function twice.
We will also need to diffrentiate between the two clients based on their computer assigned port values.
*/
func main() {
	ipAddress := flag.String("ipAddress", "Default", "Input IP Address")
	port := flag.Int("port", 0000, "Input Port Number")
    flag.Parse()
}

/*
This function will be called twice
*/
func client() {

}

/*
This function will be called once
*/
func server() {

}
