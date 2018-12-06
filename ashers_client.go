package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	allPlayers    []string
	alivePlayers  []string
	deadPlayers   []string
	myName        string = "placeholder"
	nameSet       bool   = false
	gameActive    bool   = false
	myHealth      int    = 100
	mykiller      string = "stillAlive"
	killed        int    = 0
	spectatorMode bool   = false
)

func handleGameString(str string) []byte { //handles relevant string data from messaging system
	str = strings.TrimSpace(str)
	commands := strings.Split(str, ";") //split strings by ";" separated values

	finalValue := "\n"
	switch {
	case commands[0] == "error", commands[0] == "name":
		break
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
			return []byte("error;bad args;kill\n")
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
				finalValue = fmt.Sprint("You were attacked by ", commands[2], " for ", commands[3], " damage. \n ")
			} else {
				finalValue = fmt.Sprint("You were killed by ", commands[2], "!!!\n")
				spectatorMode = true
				//TODO: report to system thatplayer died
			}
		}

	default:
		finalValue = "error;command_not_implemented\n"
		log.Print("Bad game data: ", str, "\n")
	}
	return []byte(finalValue)
}

func handleInputString(str string) []byte {
	str = strings.TrimSpace(str)
	commands := strings.Split(str, " ")

	finalValue := "\n"
	switch commands[0] {
	case "attack":
		finalValue = fmt.Sprint("attack;", commands[1], ";", myName, ";100\n")
	case "name":
		if !nameSet {
			nameSet = true
			myName = commands[1]
			finalValue = fmt.Sprint("name;", commands[1], "\n")
		} else {
			log.Print("You can only set your name once!\n")
		}
	default:
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

			if isOutgoing {
				str := string(buf[0:nBytes])
				data := handleInputString(str)
				_, err = dst.Write(data)
			} else {
				str := string(buf[0:nBytes])
				data := handleGameString(str)
				_, err = dst.Write(data)
			}
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

func StartClient(host string, port string) {
	con, err := net.Dial("tcp", host+port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to ", host+port)
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
