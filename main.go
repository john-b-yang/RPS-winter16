/*
 * Rock Paper Scissors
 * Best 2 out of 3
 */
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"flag"
	"rand"
)

/*
Server ipAddress and port can be parameters
Client port will be assigned.
For two people to play, we want to call the server function once and client function twice.
We will also need to diffrentiate between the two clients based on their computer assigned port values.
*/
func main() {
	ipAddress := flag.String("ipAddress", "169.229.50.175", "Input IP Address")
	port := flag.Int("port", 8333, "Input Port Number")
    flag.Parse()
}

func client(ipAddress string, port int) {
	numOfGames := 3
	playerScore := 0
	opponentScore := 0

	clientConn, err := net.Dial("tcp", ipAddress)
	if err != nil {
		fmt.Println("Client Connection Error: ", err)
		return
	}
	reader := bufio.NewReader(clientConn)

	//Print out connection
	fmt.Fprintf(clientConn, "GET / HTTP/1.0\r\n\r\n")
	fmt.Println("Client Connection Established Successfully, Get Ready to Play!")

	playerMove := askForPlay() //Retrieve Player Choice
	opponentMove := opponentAskForPlay()
	fmt.Println("Player picked ", playerMove, " opponent picked ", opponentMove, ". ")
	determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, numOfGames)
}

/* Client Helper Functions */
func askForPlay() {
	for {
		fmt.Println("Please type in R (Rock), P (Paper), or S (Scissors)")
		playPointer := flag.String("Play", "None", "Enter R, P, or S")
		flag.Parse()

		if *playPointer != "R" && *playPointer != "P" && *playPointer != "S" {
			fmt.Println("Your choice cannot be interpretted")
		} else {
			return *playPointer
		}
	}
}

func opponentAskForPlay() {
	var moveDictionary := map[int]string {0: "R", 1: "P", 2: "S"}
	return moveDictionary[rand.Intn(3)]
}

func determineRoundWinner(playerMove string, opponentMove string, playerScore int, opponentScore int, numOfGames int) {
	numOfGames -= 1
	if playerMove == opponentMove {
		fmt.Println("Draw! An extra game will be played!")
		numOfGames += 1
	}
	else if (playerMove == "R" && opponentMove == "S") || (playerMove == "S" && opponentMove == "P") || (playerMove == "P" && opponentMove == "R") {
		fmt.Println("Player Wins!")
		playerScore += 1
	}
	else {
		fmt.Println("Opponent Wins!")
		opponentScore += 1
	}
}

/*
This function will be called once
*/
func server() {

}
