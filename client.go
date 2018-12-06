package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	usageString string = "\n\nPossible commands:\nhelp\t\tDisplay this message.\nlist\t\tList all active players.\nattack\t\tAttack some player.\nname\t\tSet your name. Can only be called once.\nscore\t\tOutput current score of game."
)

var (
	alivePlayers  []string
	deadPlayers   []string
	allPlayers    []string
	myName        = "placeholder"
	nameSet       = false
	gameActive    = false
	myHealth      = 100
	mykiller      = "stillAlive"
	spectatorMode = false
)

func handleGameString(str string) []byte { //handles relevant string data from messaging system
	str = strings.TrimSpace(str)
	commands := strings.Split(str, ";") //split strings by ";" separated values

	finalValue := "\n"
	switch {
	case commands[0] == "name":
		allPlayers = append(allPlayers, commands[1])
		log.Println("New player added.")
	case commands[0] == "meta": //a message type for anything else
		switch {
		case commands[1] == "all players":
			copy(allPlayers, commands[2:])
			if spectatorMode {
				finalValue = fmt.Sprint("Updated ", commands[1], " to game state.\n")
			}

		case commands[1] == "alive players":
			copy(alivePlayers, commands[2:])
			if spectatorMode {
				finalValue = fmt.Sprint("Updated ", commands[1], " to game state.\n")
			}

		case commands[1] == "dead players":
			copy(deadPlayers, commands[2:])
			finalValue = fmt.Sprint("Updated ", commands[1], " to game state.\n")

		case commands[1] == "game state":
			switch commands[2] {
			case "start":
				gameActive = true
				finalValue = "Game has began!\n"

			case "end":
				gameActive = false
				finalValue = "Game over!\n"
			}

		default:
			finalValue = "error;bad meta tag\n"
		}
	case commands[0] == "start": //start game
		gameActive = true
	case commands[0] == "stop":
		if len(commands) < 2 {
			return []byte("error, should report final game state in Stop signal\n")
		}
		//finalValue = fmt.Sprint(commands[1])
		gameActive = false
	case commands[0] == "death": //reports another players death
		if len(commands) < 3 { //error check
			return []byte("error;bad args;death\n")
		}
		finalValue = fmt.Sprint("Player ", commands[2], " has been killed by ", commands[1], "\n") //output who died
	case commands[0] == "attack": //a player was attacked
		if len(commands) < 4 {
			return []byte("error;bad args;attack\n")
		}
		if spectatorMode {
			log.Println(commands[1], " was attacked by ", commands[2], " for ", commands[3], "damage.\n")
		} else if commands[1] == myName { //if this player was attacked
			damage, _ := strconv.Atoi(commands[3])
			if myHealth > damage {
				myHealth = myHealth - damage
				finalValue = fmt.Sprint("You were attacked by ", commands[2], " for ", commands[3], " damage.\n")
			} else {
				finalValue = fmt.Sprint("You were killed by ", commands[2], "!!!\n")
				spectatorMode = true
				//TODO: report to system thatplayer died
			}
		}

	default:
		finalValue = "error;command not implemented\n"
		log.Print("Bad game data: ", str, "\n")
	}
	return []byte(finalValue)
}

func handleInputString(str string) []byte {
	str = strings.TrimSpace(str)
	commands := strings.Split(str, " ")

	finalValue := ""
	switch commands[0] {
	case "help":
		log.Print(usageString)
	case "raw":
		finalValue = commands[1] + "\n"
	case "clear":
		cmd := exec.Command("clear")
		cmd.Run()
	case "attack":
		if !gameActive {
			log.Print("Game has not started yet!\n")
		} else if myHealth == 0 {
			log.Print("You are dead and cannot attack anymore!\n")
		} else if len(commands) < 3 {
			log.Print("You have not given enough arguments!\n")
			log.Print("Format: attack [PLAYER NAME] [DAMAGE]\n")
		} else {
			finalValue = fmt.Sprint("attack;", commands[1], ";", myName, ";", commands[2], "\n")
		}
	case "name":
		if len(commands) < 2 {
			log.Println("My name: ", myName)
			break
		}

		if nameSet {
			log.Print("You can only set your name once!\n")
		} else if len(commands[1]) > 5 {
			log.Print("Your name must be 5 characters or less.\n")
		} else {
			found := false
			for _, v := range allPlayers {
				if v == commands[1] {
					log.Println("Another player has already taken that name!")
					found = true
				}
			}
			if found {
				break
			}
			nameSet = true
			myName = commands[1]
			finalValue = fmt.Sprint("name;", commands[1], "\n")
			log.Println("Name has successfully been set.")
		}
	case "list":
		if len(allPlayers) == 0 {
			log.Println("No other players have registered on this server yet.")
		} else {
			log.Print("My name: ", myName, "\n")
			log.Print("All names:\n")
			for i, v := range allPlayers {
				log.Print("Name ", i, ": ", v)
			}
		}
	case "score":
		log.Println("Score not implemented yet...")
	default:
		log.Print("Error: Bad input.\n")
		log.Println(usageString)
		finalValue = fmt.Sprint("error;bad_input_string;", str, "\n")
	}

	return []byte(finalValue)
}

func streamCpy(src io.Reader, dst io.Writer, isOutgoing bool) <-chan int {
	buf := make([]byte, 1024)
	sync := make(chan int)

	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				log.Printf("Con from %v closed\n", con.RemoteAddr())
			}
			sync <- 0
		}()

		for {

			nBytes, err := src.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Read error: %s\n", err)
				}
				break
			}

			var data []byte
			if isOutgoing {
				str := string(buf[0:nBytes])
				data = handleInputString(str)
			} else {
				str := string(buf[0:nBytes])
				data = handleGameString(str)
			}

			_, err = dst.Write(data)

			if err != nil {
				log.Fatalf("Write error: %s\n", err)
			}
		}
	}()

	return sync

}

func HandleCons(con net.Conn) {

	stdoutChan := streamCpy(con, os.Stdout, false)
	remoteChan := streamCpy(os.Stdin, con, true)

	select {
	case <-stdoutChan:
		log.Println("Remote connection broken.")
	case <-remoteChan:
		log.Println("Local connection broken.")
	}
}

// spawns client goroutines that coordinate streams
func StartClient(host string, port string) {
	con, err := net.Dial("tcp", host+port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to ", host+port)
	log.Println("Welcome to GoWar!!", usageString)
	allPlayers = append(allPlayers, "admin")
	HandleCons(con)
}

func main() {

	var host, port string
	flag.StringVar(&host, "host", "", "Remote host to connect to")
	flag.StringVar(&port, "port", ":8080", "Port of remote host")
	flag.Parse()

	if host == "" {
		flag.Usage()
	}

	StartClient(host, port)

}
