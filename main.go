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
	//Playing against Ashwarya
	JohnIPAddress := flag.String("John's ipAddress", "169.229.50.175", "Input IP Address")
	JohnPort := flag.Int("John's port", 6421, "Input Port Number")
	AshIPAddress := flag.String("Ash's ipAddress", "169.229.50.188", "Input IP Address")
	AshPort := flag.String("Ash's port", 8333, "Input Port Number")
    flag.Parse()
	//Organize this into conditionals

	//If I play against myself
	client(*JohnIPAddress, *JohnPort)
	//If I play against Ashwarya
	client(*AshIPAddress, *AshPort)
	//Instantiating Server code
	server(*JohnPort)
}

/*
So this client currently only works between a human and computer
But we don't really want that
*/
func client(ipAddress string, port int) {
	//Concatenating ipAddress and port number
	iPAddPort := fmt.Sprintf("%s:%d", ipAddress, port)
	//Create and test client connection
	clientConn, err := net.Dial("tcp", iPAddPort)
	if err != nil {
		fmt.Println("Client Connection Error: ", err)
		return
	} else {
		fmt.Println("Client Connection Established")
	}
	reader := bufio.NewReader(clientConn)
	//Print out connection
	fmt.Fprintf(clientConn, "GET / HTTP/1.0\r\n\r\n")
	fmt.Println("Client Connection Established Successfully, Get Ready to Play!")

	numOfGames := 3
	playerScore := 0
	opponentScore := 0

    for round := 0; round < numOfGames; round++ {}
		playerMove := askForPlay() //Retrieve Player Choice
		opponentMove := opponentAskForPlay()
		fmt.Println("Player picked ", playerMove, " opponent picked ", opponentMove, ". ")
		determineRoundWinner(playerMove, opponentMove, playerScore, opponentScore, round)
		printStage(playerScore, opponentScore)

		if _, err := clientConn.Write([]byte(playerMove)); err != nil {
	        fmt.Println("Send failed:", err)
	        os.Exit(1)
	    }
	}
	clientConn.close()
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

func determineRoundWinner(playerMove string, opponentMove string, playerScore int, opponentScore int, round int) {
	round -= 1
	if playerMove == opponentMove {
		fmt.Println("Draw! An extra game will be played!")
		round += 1
	}
	else if playerMove == "R" && opponentMove == "S", playerMove == "S" && opponentMove == "P", playerMove == "P" && opponentMove == "R" {
		fmt.Println("Player wins this round!")
		playerScore += 1
	}
	else {
		fmt.Println("Opponent wins this round!")
		opponentScore += 1
	}
}

func printStage(playerScore int, opponentScore int) ->  {
	if playerScore == 2 {
		fmt.Printf("Player wins the game by a score of (%d)-(%d)!", playerScore, opponentScore)
	} else if opponentScore == 2 {
		fmt.Printf("Opponent wins the game by a score of (%d)-(%d)!", opponentScore, playerScore)
		return true
	} else {
		fmt.Println("Next Round!")
		return false
	}
}

/*
This function will be called once
*/
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
		}
		fmt.Printf("(%d) Received: %s", i, string(recvMsgBytes))

		//Sending Message, to be modified for RPS
		message := string(recvMsgBytes)
		if message == nil {
			sengMsg := "Nil Message\n"
		} else if message == "R" {
			sengMsg := "P\n"
		} else if message == "P" {
			sengMsg := "S\n"
		} else if message == "S" {
			sengMsg := "Rt\n"
		}
		fmt.Printf("(%d) Sending: %s\n", i, sengMsg) //MARK

		if _, err := serverConn.Write([]byte(sendMsg)); err != nil {
			fmt.Println("Send failed:", err)
			os.Exit(1)
		}
	}
	serverConn.close()
}
