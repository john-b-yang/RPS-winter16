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
	"flag"
	"math/rand"
	"time"
	"strings"
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
	gameMode := flag.String("gameMode", "filler", "Enter Player or CPU")
	player := flag.String("player", "filler", "Please indicate whether you are a Client or Server")
	opponent := flag.String("opponent", "filler", "Is your opponent a Client or Server?")
    port := flag.Int("port", 6421, "6421 (John) or 8333 (Ashwarya)")
    ipAddress := flag.String("ipAddress", "169.229.50.175", "169.229.50.175 (John) or 169.229.50.188 (Ashwarya)")

	flag.Parse()

	//Logic for calling client/server
	if *player != "filler" && *gameMode != "filler" && *opponent != "filler" {
		if *gameMode == "player" { //Person opponent
			if *player == "server" {
				fmt.Println("Starting Server in Interactive Mode")
				serverPlay(*port) //Ash's Port
			} else if *player == "client" {
				fmt.Println("Starting Client in Interactive Mode")
				clientPlay(*ipAddress, *port) //Ash's IP Address + Port
			} else {
				fmt.Println("Please enter an appropriate player type")
			}
		} else if *gameMode == "CPU" { //CPU opponent
			if *player == "server" {
				fmt.Println("Starting Server in Automatic Mode")
				serverCPU(*port) //John's Port
			} else if *player == "client" {
				fmt.Println("Starting Client in Automatic Mode")
				clientCPU(*ipAddress, *port) //John's IP Address + Port
			} else {
				fmt.Println("Please enter an appropriate player type")
			}
		} else {
			fmt.Println("Please choose an appropriate game mode")
		}
	} else { //Error Message
		fmt.Println("Please answer all fields for the game to begin correctly")
	}
}

func clientCPU(ipAddress string, port int) {
	ipAddPort := fmt.Sprintf("%s:%d", ipAddress, port) //Concatenating ipAddress and port number
	clientConn, err := net.Dial("tcp", ipAddPort) //Create and test client connection
	if err != nil {
		fmt.Println("Client Connection Error: ", err)
		return
	}
	fmt.Println("Client Connection Established Successfully, Get Ready to Play!")

	numOfGames := 3
	playerScore := 0
	opponentScore := 0

	fmt.Println("Rock, Paper, Scissors!")
	fmt.Println("----------------------")
	//Figure out how to terminate this loop
    for round := 0; round < numOfGames; round++ {
		playerMove := opponentAskForPlay() //Retrieve Player Choice

		if _, err := clientConn.Write([]byte(playerMove + "\n")); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}

		reader := bufio.NewReader(clientConn)
		//Receiving Message
		recvMsgBytes, err := reader.ReadBytes('\n') //recvMsg is opponent play

		if err != nil {
			fmt.Println("Error reading opponent play: ", err)
			return
		}

		opponentMove := string(recvMsgBytes)
		opponentPrint := strings.TrimSpace(opponentMove)
		printRound := round + 1
		fmt.Println("Game", printRound, ": Player played", playerMove, "and Opponent played", opponentPrint)

		result := determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, round) //Increment round number accordingly

		if result == "tie" {
			round -= 1
		} else if result == playerMove {
			playerScore += 1
		} else if result == opponentMove {
			opponentScore += 1
		}

		printStage(playerScore, opponentScore) //Checks whether one of the players has won
		fmt.Println("----------------------")
		time.Sleep(2 * time.Second)
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

	reader := bufio.NewReader(serverConn)
	numOfGames := 3
	playerScore := 0
	opponentScore := 0

	fmt.Println("Rock, Paper, Scissors!")
	fmt.Println("----------------------")
	for i := 0; i < numOfGames; i++ {
		//Received Message
		recvMsgBytes, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Receive failed:", err)
			os.Exit(1)
		}

		opponentMove := string(recvMsgBytes)
		opponentPrint := strings.TrimSpace(opponentMove)
		playerMove := opponentAskForPlay2(opponentMove)

		//Sending Message
		if _, err := serverConn.Write([]byte(playerMove + "\n")); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}
		printRound := i + 1
		fmt.Println("Game", printRound, ": Player played", playerMove, "and Opponent played", opponentPrint)
		result := determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, i)

		if result == "tie" {
			i -= 1
		} else if result == playerMove {
			playerScore += 1
		} else if result == opponentMove {
			opponentScore += 1
		}

		printStage(playerScore, opponentScore)
		fmt.Println("----------------------")
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

func opponentAskForPlay2(opponentMove string) string {
	if opponentMove == "R\n" {
		return "S"
	} else if opponentMove == "P\n" {
		return "R"
	} else if opponentMove == "S\n" {
		return "P"
	} else {
		return opponentAskForPlay()
	}
}

//Determine the winner of a given round
//MARK: Definitive issue with this method, comparisons not correct
func determineRoundWinner(playerMove string, opponentMove string, playerScore int, opponentScore int, round int) string {
	// fmt.Println("Determine Round Winner Entered")
	// fmt.Println("s", playerMove, "s")
	// fmt.Println("s", opponentMove, "s")
	// check := (playerMove + "\n") == opponentMove
	// fmt.Println(check)
	if playerMove + "\n" == opponentMove {
		fmt.Println("Draw! An extra game will be played!")
		return "tie"
	} else if (playerMove == "R\n" && opponentMove == "S") || (playerMove == "S\n" && opponentMove == "P") || (playerMove == "P\n" && opponentMove == "R") {
		fmt.Println("Player wins this round!")
		return playerMove
	} else {
		fmt.Println("Opponent wins this round!")
		return opponentMove
	}
}

func printStage(playerScore int, opponentScore int) {
	if playerScore == 2 {
		fmt.Println("Player wins the game by a score of", playerScore, "-", opponentScore, "!")
		os.Exit(1)
	} else if opponentScore == 2 {
		fmt.Println("Opponent wins the game by a score of", opponentScore, "-", playerScore, "!")
		os.Exit(1)
	} else {
		fmt.Println("Next Round! Current score of player vs opponent is", playerScore, "-", opponentScore, "!")
	}
}

func clientPlay(ipAddress string, port int) {
    print("Filler")
}

func serverPlay(port int) {
    print("Filler")
}
