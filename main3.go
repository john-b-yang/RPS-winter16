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

/*
	Things to work on
	1. Network connection issues
	2. Randomness: seed?
	3. Algorithm
	4. Accepting/Sending different types of input between Ashwarya's code and mine
*/

/* Command Line Prompts:
	For CPU/Automatic Mode:
		-Client: ./main2 -gameMode=CPU -port=6421 -ipAddress=127.0.0.1 -player=client
		-Server: ./main2 -gameMode=CPU -port=6421 -ipAddress=127.0.0.1 -player=server
	For Player/Interactive Mode:
		-Client: ./main2 -gameMode=player -port=6421 -ipAddress=127.0.0.1 -player=client
		-Server: ./main2 -gameMode=player -port=6421 -ipAddress=127.0.0.1 -player=server
*/
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
    port := flag.Int("port", 6421, "6421 (John) or 8333 (Ashwarya)")
    ipAddress := flag.String("ipAddress", "127.0.0.1", "169.229.50.175 (John) or 169.229.50.188 (Ashwarya)")
	flag.Parse()

	//Logic for calling client/server
	if *player != "filler" && *gameMode != "filler" {
		if *gameMode == "player" { //Person opponent
			if *player == "server" {
				fmt.Println("Starting Server in Interactive Mode")
				server(*port, false) //Ash's Port
			} else if *player == "client" {
				fmt.Println("Starting Client in Interactive Mode")
				client(*ipAddress, *port, false) //Ash's IP Address + Port
			} else {
				fmt.Println("Please enter an appropriate player type")
			}
		} else if *gameMode == "CPU" { //CPU opponent
			/*
				In automatic mode, the way it is currently configured, the server will always win. Originally, I wrote the code so that
				both client and server would call the opponentAskForPlay() function which would generate a random number => random move
				However, since the ipAddress for the computers were the same, and computers are pseudo-random (time factor), they kept generating
				the same numbers, and therefore, the same moves, so instead, I wrote a function opponentAskForPlay2() that would generate a move with
				knowledge of the opponent's move.

				To make the automatic mode more random, I would have to
				1. Use different ipAddresses in the client and server calls
				2. Change the call to "opponentAskForPlay2()" in the server method to just "opponentAskForPlay()"
			*/
			if *player == "server" {
				fmt.Println("Starting Server in Automatic Mode")
				server(*port, true) //John's Port
			} else if *player == "client" {
				fmt.Println("Starting Client in Automatic Mode")
				client(*ipAddress, *port, true) //John's IP Address + Port
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

func client(ipAddress string, port int, isCPU bool) {
	ipAddPort := fmt.Sprintf("%s:%d", ipAddress, port) //Concatenating ipAddress and port number
	clientConn, err := net.Dial("tcp", ipAddPort) //Create and test client connection
	if err != nil {
		fmt.Println("Client Connection Error: ", err)
		return
	}
	fmt.Println("Client Connection Established Successfully, Get Ready to Play!\n")

	numOfGames := 3
	playerScore := 0
	opponentScore := 0

	fmt.Println("Rock, Paper, Scissors!")
	fmt.Println("----------------------")
	for round := 0; round < numOfGames; round++ {
		var playerMove string = "filler"
		if isCPU {
			playerMove = opponentAskForPlay()
		} else {
			playerMove = askForPlay()
		}

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
		} else if result == "error" {
			fmt.Println("Technical Difficulty")
			os.Exit(1)
		}

		printStage(playerScore, opponentScore) //Checks whether one of the players has won
		fmt.Println("----------------------")
		time.Sleep(2 * time.Second)
	}
	clientConn.Close()
}

//Server Function
func server(port int, isCPU bool) {
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
	fmt.Println("Server Connection Established Successfully, Get Ready to Play!\n")

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

		//Retrieve Opponent Move and generate Player Move
		opponentMove := string(recvMsgBytes)
		var playerMove string = "filler"
		if isCPU {
			playerMove = opponentAskForPlay2(opponentMove)
		} else {
			playerMove = askForPlay()
		}
		//Sending Message
		if _, err := serverConn.Write([]byte(playerMove + "\n")); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}

		//Print Formatting
		opponentPrint := strings.TrimSpace(opponentMove)
		printRound := i + 1
		fmt.Println("Game", printRound, ": Player played", playerMove, "and Opponent played", opponentPrint)
		result := determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, i)

		if result == "tie" {
			i -= 1
		} else if result == playerMove {
			playerScore += 1
		} else if result == opponentMove {
			opponentScore += 1
		} else if result == "error" {
			os.Exit(1)
		}

		printStage(playerScore, opponentScore)
		fmt.Println("----------------------")
	}
	serverConn.Close()
}

/********************/
/* Helper Functions */
/********************/

//Prompt user for play
func askForPlay() string {
	fmt.Println("Please type in R (Rock), P (Paper), or S (Scissors)")
	var play string
	fmt.Scanln(&play)
	return play
}

//Automatic Opponent (For Automatic Mode)
func opponentAskForPlay() string {
	moveDictionary := map[int]string {0: "R", 1: "P", 2: "S"}
	return moveDictionary[rand.Intn(3)]
}

//Automatic Opponent V2
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
func determineRoundWinner(playerMove string, opponentMove string, playerScore int, opponentScore int, round int) string {
	if playerMove + "\n" == opponentMove {
		fmt.Println("Draw! An extra game will be played!")
		return "tie"
	} else if (playerMove == "R" && opponentMove == "S\n") || (playerMove == "S" && opponentMove == "P\n") || (playerMove == "P" && opponentMove == "R\n") {
		fmt.Println("Player wins this round!")
		return playerMove
	} else if (playerMove == "S" && opponentMove == "R\n") || (playerMove == "P" && opponentMove == "S\n") || (playerMove == "R" && opponentMove == "P\n") {
		fmt.Println("Opponent wins this round!")
		return opponentMove
	} else if playerMove == "exit" {
		fmt.Println("Exiting game!")
		return "error"
	} else {
		fmt.Println("Input was not understood. Game is terminating.")
		return "error"
	}
}

//Print out the point in game flow after a round is complete
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
