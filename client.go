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
	allPlayers       []string
	playerPoints     []int
	alivePlayers     []string
	myName           = "placeholder"
	nameSet          = false
	gameActive       = false
	myHealth         = 100
	mykiller         = "stillAlive"
	spectatorMode    = false
	sendDeathMessage = false
)

func isPlayerAlive(name string) bool { //return true if named player is in alivePlayers list, false if not
	for i := 0; i < len(alivePlayers); i++ {
		if name == alivePlayers[i] {
			return true
		}
	}
	return false
}

func printAliveList() {
	log.Println("Current alive players: \n")
	for i := range alivePlayers {
		log.Println(alivePlayers[i])
	}
}

func killPlayer(name string) bool { // takes out player from alive Players list
	for i := 0; i < len(alivePlayers); i++ {
		if name == alivePlayers[i] {
			alivePlayers = append(alivePlayers[:i], alivePlayers[i+1:]...) //take out player from list
			return true
		}
	}
	log.Println("Error! player", name, "not in the alive players list\n")
	return false
}

func playerScore(name string) int { //returns score for given player
	if len(allPlayers) != len(playerPoints) {
		log.Println("ERROR players and points arrays not same size\n")
		return 0
	}
	for i := 0; i < len(allPlayers); i++ {
		if myName == allPlayers[i] {
			return playerPoints[i]
		}
	}
	log.Println("Error! You're not in your own list!!!\n")
	return 0
}

func givePoints(name string, points int) { //give points to specific player
	if len(allPlayers) != len(playerPoints) {
		log.Println("ERROR players and points arrays not same size\n")
		return
	}
	for i := 0; i < len(allPlayers); i++ {
		if name == allPlayers[i] {
			playerPoints[i] = playerPoints[i] + points
			return
		}
	}
	log.Println("Error in givePoints(),  player", name, "not in the players list\n")
}

func outPutScores() string {

	if len(allPlayers) != len(playerPoints) {
		log.Println("ERROR players and points arrays not same size\n")
		return "err\n"
	}
	message := "GAME OVER!!\n"
	if len(alivePlayers) == 1 {
		message += "Last living player: " + alivePlayers[0]
	} else if len(alivePlayers) > 1 {
		message += "Players left standing:"
		for i := 0; i < len(alivePlayers); i++ {
			message += alivePlayers[i] + ", "
		}
	} else {
		message += "NO players left standing!"
	}
	message += "\n Score Board \n ----------------------------\n"
	for i := 0; i < len(allPlayers); i++ {
		message += allPlayers[i] + "--------" + strconv.Itoa(playerPoints[i]) + "\n"
	}
	return message
}

//handles information given from public record and handles relevate data
func handleGameString(str string) []byte {
	str = strings.TrimSpace(str)
	commands := strings.Split(str, ";") //split strings by ";" separated values

	finalValue := ""
	switch {
	case commands[0] == "name":
		allPlayers = append(allPlayers, commands[1])
		playerPoints = append(playerPoints, 0)
		log.Println("New player added.")
	case commands[0] == "meta": //a message type for anything else
		if len(commands) == 2 {
			finalValue = fmt.Sprint(commands[1], "\n") //if meta should be read as just a message, then print message
		} else if len(commands) > 2 {
			if commands[1] == "all players" { //if meta is an update of all
				allPlayers = commands[2:]//recreate list of players
				playerPoints :=make([]int,len(allPlayers))//clears points then recreate in loop below
				for i:= range allPlayers { //make a list of player points all starting off at 0
					playerPoints[i] = 100
				}
				alivePlayers = allPlayers
				if spectatorMode {
					finalValue = fmt.Sprint("Updated players.\n")
				}
			}
		}
	case commands[0] == "start": //start game
		finalValue = fmt.Sprint("GoWar STARTED!!! \n")
		gameActive = true
	case commands[0] == "stop"://stop game and output scores
		if gameActive {
			finalValue = fmt.Sprint(outPutScores()) //print game results given in stop command
			gameActive = false
		}

	case commands[0] == "death": //reports another players death
		if len(commands) < 3 { //error check
			return []byte("error;bad args;death\n")
		}
		finalValue = fmt.Sprint("Player ", commands[1], " has been killed by ", commands[2], "\n") //output who died
		killPlayer(commands[1])
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
				mykiller = commands[2]
				sendDeathMessage = true
				myHealth = 0
				spectatorMode = true
				//TODO: report to system that player died
			}
		}
	default:
		// finalValue = "error;unhandled tag;" + commands[0] + "\n"
		finalValue = ""
	}
	return []byte(finalValue)
}

//handles messages given from this client, and processes them to be sent to other players
func handleInputString(str string) []byte {
	str = strings.TrimSpace(str)
	commands := strings.Split(str, " ")
	if sendDeathMessage { //if they died then first loop will send death message once
		return []byte(fmt.Sprint("death;", myName, ";", mykiller, "\n"))
		sendDeathMessage = false

	}
	if myHealth == 0 { //check first to see if player is still alive
		log.Print("You are dead :(\n")
		return []byte(fmt.Sprint("\n"))
	}
	finalValue := "\n"
	switch commands[0] {
	case "help": //print list of usable commands
		log.Print(usageString)
	case "raw": //enter custom data without user error thrown, hackers here ya go :)
		finalValue = commands[1] + "\n"
	case "clear": //clear console
		cmd := exec.Command("clear")
		cmd.Run()
	case "attack": // attack a player, if bad format, or not alive, give error message
		if !gameActive {
			log.Print("Game has not started yet!\n")
		} else if len(commands) < 2 {
			log.Print("You have not given enough arguments!\n")
			log.Print("Format: attack [PLAYER NAME]\n")
		} else if isPlayerAlive(commands[1]) {
			finalValue = fmt.Sprint("attack;", commands[1], ";", myName, "; 10\n")
		} else {
			log.Print(commands[1], "is dead and gone.")
		}
	case "name": // set name if given proper input
		if len(commands) < 2 {
			log.Println("My name: ", myName)
			break
		}

		if nameSet {
			log.Print("You can only set your name once!\n")
		} else if len(commands[1]) > 5 {
			log.Print("Your name must be 5 characters or less with no special characters.\n")
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
			allPlayers = append(allPlayers, myName)
			playerPoints = append(playerPoints, 0)
			finalValue = fmt.Sprint("name;", commands[1], "\n")
			log.Println("Name has successfully been set.")
		}
	case "list": //lists all players in the server
		if len(allPlayers) == 0 {
			log.Println("No other players have registered on this server yet.")
		} else {
			log.Print("Players:\n")
			for i := 0; i < len(allPlayers); i++ {
				log.Print("Player ", i, ": ", allPlayers[i])
			}
		}
	case "score": //return players score
		log.Println("My score: ", strconv.Itoa(playerScore(myName)), "\n")
	case "spec": //enter spectator mode
		myHealth = 0
		spectatorMode = true
	default: //return error if bad string given
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
