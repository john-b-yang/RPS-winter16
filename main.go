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

TBD List
-IP address
-Port Number
*/
func client() {
	clientConn, err := net.Dial("tcp", "TBD: IP ADDRESS")
	if err != nil {
		fmt.Println("Client Connection Error: ", err)
		return
	}
	//Print out connection
	fmt.Fprintf(clientConn, "GET / HTTP/1.0\r\n\r\n")
	fmt.Println("Client Connection Established Successfully, Get Ready to Play!")

	play := askForPlay() //Retrieve Player Choice

	/* TBD: The code for diffrentiating which client is called from which computer. Also must establish how we compare one player's code to another

	Another matter, do we print different messages to different terminals? For example, if player picks rock and opponent picks scissors, do the two player see the same message or different outputs?*/

	//FOR NOW we will use "opponent" to represent other player's pick.
	opponent := "Filler"
	fmt.Println("Player picked ", play, " opponent picked ", opponent, ". ")
	if play == opponent {
		fmt.Println("Draw!")
	}
	else if (play == "R" && opponent == "S") || (play == "S" && opponent == "P") || (play == "P" && opponent == "R") {
		fmt.Println("Player Wins!")
	}
	else {
		fmt.Println("Opponent Wins!")
	}
}

/* Client Helper Functions */
func askForPlay() {
	for {
		fmt.Println("Please type in R (Rock), P (Paper), or S (Scissors)")
		playPointer := flag.String("Play", "None", "Enter R, P, or S")
		flag.Parse()

		if *playPointer != "R" && play != "P" && play != "S" {
			fmt.Println("Your choice cannot be interpretted")
		} else {
			return *playPointer
		}
	}
}

/*
This function will be called once
*/
func server() {

}
