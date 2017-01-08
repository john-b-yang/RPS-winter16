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

func main() {
	/*
	John's IP Address: 169.229.50.175
	John's Port: 6421
	Ashwarya's IP Address: 169.229.50.188
	Ashwarya'as IP Address: 8333
	*/

	//Diffrentiate between which port, ipAddress to use
	gameMode := flag.String("Game Mode", "Player", "Enter Player or CPU")
	//Port + IP Address: John's then Ashwarya's
	port := flag.Int("Port", 64218333, "6421 (John) or 8333 (Ashwarya)")
	ipAddress := flag.String("IP Address", "", "169.229.50.175 (John) or 169.229.50.188 (Ashwarya)")
	//Player Type
	player := flag.String("Player", "temp", "Please indicate whether you are a Client or Server")
	opponent := flag.String("opponent", "temp", "Is your opponent a Client or Server?")

	flag.Parse()

	//Logic for calling client/server
	if *player != "" && *gameMode != "" && *opponent != "" {
		if *gameMode == "Player" { //Person opponent
			if *player == "server" {
				server(*port) //Ash's Port
			} else if *player == "client" {
				client(*ipAddress, *port) //Ash's IP Address + Port
			}
		} else if *gameMode == "CPU" { //CPU opponent
			if *player == "server" {
				server(*port) //John's Port
			} else if *player == "client" {
				client(*ipAddress, *port) //John's IP Address + Port
			}
		}
	} else { //Error Message
		fmt.Println("Please answer all fields for the game to begin correctly")
		return
	}
}

func client(ipAddress string, port int) {
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

	reader := bufio.NewReader(clientConn)
	numOfGames := 3 //Should this be 2? numOfGames = # of games to be won OR most # of games?
	playerScore := 0
	opponentScore := 0

	//Figure out how to terminate this loop
    for round := 0; round < numOfGames; round++ {
		//Receiving Message
		recvMsg, err := reader.ReadString('\n') //recvMsg is opponent play
		if err != nil {
			fmt.Println("Error reading next play: ", err)
			return
		}

		//Game Logic
		/*********/
		print(recvMsg) //MARK: Using recvMsg so go build will not error
		playerMove := askForPlay() //Retrieve Player Choice
		opponentMove := opponentAskForPlay()
		fmt.Println("Player picked ", playerMove, " opponent picked ", opponentMove, ". ")

		determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, round) //Increment round number accordingly
		isGameOver := printStage(playerScore, opponentScore) //Checks whether one of the players has won

		//MARK: Is this necessary + How to check if game is over
		if isGameOver {
			print("TBD: Game Over Message")
		}
		/*********/

		//Sending Message
		if _, err := clientConn.Write([]byte(playerMove)); err != nil {
	        fmt.Println("Send failed:", err)
	        os.Exit(1)
	    }
	}
	clientConn.Close()
}

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

//Automatic Opponent (So not a client)
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
func printStage(playerScore int, opponentScore int) bool  {
	if playerScore == 2 {
		fmt.Printf("Player wins the game by a score of (%d)-(%d)!", playerScore, opponentScore)
		return true
	} else if opponentScore == 2 {
		fmt.Printf("Opponent wins the game by a score of (%d)-(%d)!", opponentScore, playerScore)
		return true
	} else {
		fmt.Println("Next Round!")
		return false
	}
}

func server(port int) {
	portString := fmt.Sprintf(":%d", port)

	//Listening
	ln, err := net.Listen("tcp", portString) //Same Port Number as Client's
	if err != nil {
		fmt.Println("Listen failed:", err)
		os.Exit(1)
	} else {
		fmt.Println("Listening Passed")
	}

	//Accepting
	serverConn, err := ln.Accept()
	if err != nil {
		fmt.Println("Accept failed:", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(serverConn)

	numOfGames := 3
	for i := 0; i < numOfGames; i++ {
		//Received Message
		recvMsgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Receive failed:", err)
			os.Exit(1)
		} else {
			fmt.Printf("(%d) Received: %s", i, string(recvMsgBytes))
		}

		//Sending Message, to be modified for RPS
		message := string(recvMsgBytes)
		var sendMsg string = "Nil Message\n"
		if message == "R" {
			sendMsg = "P\n"
		} else if message == "P" {
			sendMsg = "S\n"
		} else if message == "S" {
			sendMsg = "R\n"
		}

		//Sending Message
		fmt.Printf("(%d) Sending: %s\n", i, sendMsg) //MARK
		if _, err := serverConn.Write([]byte(sendMsg)); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}
	}
	serverConn.Close()
}
