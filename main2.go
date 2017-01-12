/*
 * Rock Paper Scissors
 * Best 2 out of 3
 */
package main

/* Unused Packages
"time"
*/
import (
	"bufio"
	"fmt"
	"net"
	"os"
	"flag"
	"math/rand"
)

//Command Line Prompt: ./main -gameMode=CPU -port=6421 -ipAddress=169.229.50.175 -player=client

func main() {
	/*
	John's IP Address: 169.229.50.175
	John's Port: 6421
	Ashwarya's IP Address: 169.229.50.188
	Ashwarya'as IP Address: 8333
	*/

	//Parameters
	gameMode := flag.String("gameMode", "Player", "Enter Player or CPU")
	player := flag.String("player", "temp", "Please indicate whether you are a Client or Server")
	opponent := flag.String("opponent", "temp", "Is your opponent a Client or Server?")
    port := flag.Int("port", 64218333, "6421 (John) or 8333 (Ashwarya)")
    ipAddress := flag.String("ipAddress", "", "169.229.50.175 (John) or 169.229.50.188 (Ashwarya)")

	flag.Parse()

	//Logic for calling client/server
	if *player != "" && *gameMode != "" && *opponent != "" {
		if *gameMode == "Player" { //Person opponent
			if *player == "server" {
				fmt.Println("Starting Server in Interactive Mode")
				serverPlay(*port) //Ash's Port
			} else if *player == "client" {
				fmt.Println("Starting Client in Interactive Mode")
				clientPlay(*ipAddress, *port) //Ash's IP Address + Port
			}
		} else if *gameMode == "CPU" { //CPU opponent
			if *player == "server" {
				fmt.Println("Starting Server in Automatic Mode")
				serverCPU(*port) //John's Port
			} else if *player == "client" {
				fmt.Println("Starting Client in Automatic Mode")
				clientCPU(*ipAddress, *port) //John's IP Address + Port
			}
		}
	} else { //Error Message
		fmt.Println("Please answer all fields for the game to begin correctly")
		return
	}
}

func clientCPU(ipAddress string, port int) {
	iPAddPort := fmt.Sprintf("%s:%d", ipAddress, port) //Concatenating ipAddress and port number
	clientConn, err := net.Dial("tcp", iPAddPort) //Create and test client connection
	if err != nil {
		fmt.Println("Client Connection Error: ", err)
		return
	} else {
		//Print out connection
		fmt.Fprintf(clientConn, "GET / HTTP/1.0\r\n\r\n")
		fmt.Println("Client Connection Established Successfully, Get Ready to Play!")
	}

	tempMove := opponentAskForPlay()
	if _, err := clientConn.Write([]byte(tempMove + "\n")); err != nil {
		fmt.Println("Send failed:", err)
		os.Exit(1)
	}

	numOfGames := 3 //Should this be 2? numOfGames = # of games to be won OR most # of games?
	playerScore := 0
	opponentScore := 0

	//Figure out how to terminate this loop
    for round := 0; round < numOfGames; round++ {
		reader := bufio.NewReader(clientConn)
		//Receiving Message
		recvMsgBytes, err := reader.ReadBytes('\n') //recvMsg is opponent play
		if err != nil {
			fmt.Println("Error reading opponent play: ", err)
			return
		}

		//Game Logic
		playerMove := askForPlay() //Retrieve Player Choice
		opponentMove := string(recvMsgBytes)
		if _, err := clientConn.Write([]byte(playerMove + "\n")); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}
		fmt.Println("Player picked ", playerMove, " opponent picked ", opponentMove, ". ")

		determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, round) //Increment round number accordingly
		printStage(playerScore, opponentScore) //Checks whether one of the players has won

		//Sending Message
	}
	clientConn.Close()
}

func serverCPU(port int) {
	portString := fmt.Sprintf(":%d", port)

	//Listening
	ln, err := net.Listen("tcp", portString)
	if err != nil {
		fmt.Println("Listen failed:", err)
		os.Exit(1)
	}

	//Accepting
	serverConn, err := ln.Accept()
	if err != nil {
		fmt.Println("Accept failed:", err)
		os.Exit(1)
	}

	tempMove := opponentAskForPlay()
	if _, err := serverConn.Write([]byte(tempMove + "\n")); err != nil {
		fmt.Println("Send failed:", err)
		os.Exit(1)
	}

	numOfGames := 3
	playerScore := 0
	opponentScore := 0

	for i := 0; i < numOfGames; i++ {
		reader := bufio.NewReader(serverConn)
		//Received Message
		recvMsgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Receive failed:", err)
			os.Exit(1)
		}

		opponentMove := string(recvMsgBytes)
		playerMove := askForPlay()

		//Sending Message
		if _, err := serverConn.Write([]byte(playerMove + "\n")); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}

		fmt.Println("Player picked ", playerMove, " opponent picked ", opponentMove, ". ")
		determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, i)
		printStage(playerScore, opponentScore)
	}
	serverConn.Close()
}

/***************************/
/* Client Helper Functions */

//Prompt user for play
func askForPlay() string {
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

//Automatic Opponent (For Automatic Mode)
func opponentAskForPlay() string {
	moveDictionary := map[int]string {0: "R", 1: "P", 2: "S"}
	return moveDictionary[rand.Intn(3)]
}

//Determine the winner of a given round
func determineRoundWinner(playerMove string, opponentMove string, playerScore int, opponentScore int, round int) {
	round += 1
	if playerMove == opponentMove {
		fmt.Println("Draw! An extra game will be played!")
		round -= 1
	} else if (playerMove == "R" && opponentMove == "S") || (playerMove == "S" && opponentMove == "P") || (playerMove == "P" && opponentMove == "R") {
		fmt.Println("Player wins this round!")
		playerScore += 1
	} else {
		fmt.Println("Opponent wins this round!")
		opponentScore += 1
	}
}

//
func printStage(playerScore int, opponentScore int) {
	if playerScore == 2 {
		fmt.Printf("Player wins the game by a score of (%d)-(%d)!", playerScore, opponentScore)
		os.Exit(1)
	} else if opponentScore == 2 {
		fmt.Printf("Opponent wins the game by a score of (%d)-(%d)!", opponentScore, playerScore)
		os.Exit(1)
	} else {
		fmt.Println("Next Round! Current score of player vs opponent is (%d)-(%d)!", playerScore, opponentScore)
	}
}

func clientPlay(ipAddress string, port int) {
    print("Filler")
}

func serverPlay(port int) {
    print("Filler")
}
