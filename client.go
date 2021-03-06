package main

import ( //import all needed libraries
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const ( //a constant string that shows the possible commands (or at least the ones we want the user to know about )
	usageString string = "\n\nPossible commands:\nhelp\t\tDisplay this message.\nlist\t\tList all active players.\nattack\t\tAttack some player.\nname\t\tSet your name. Can only be called once.\nscore\t\tOutput current score of game."
)

var ( //player specific variables
	allPlayers       []string        //list of all players in game
	playerPoints     int             //amount of points player has earned
	myName           = "placeholder" //players name, start just as placeholder
	nameSet          = false         //flag if name has been set
	gameActive       = false         //flag if game has started
	myHealth         = 100           //each player starts with 100 health points
	mykiller         = "stillAlive"  //keeps track of killer to let you know who killed ya
	spectatorMode    = false         //spectator mode lets dead players watch everyone else attack each other
	sendDeathMessage = false         //flag to report to system that user has been killed
)

func isPlayerAlive(name string) bool { //return true if named player is in alivePlayers list, false if not
	for i := 0; i < len(allPlayers); i++ { //look through whole alive list to see if given player is in there
		if name == allPlayers[i] {
			return true
		}
	}
	return false //player was not found in list
}

////print alive players function not needed
// func printAliveList() {//print off list of all alive players
// 	log.Println("Current alive players: \n")
// 	for i := range alivePlayers {
// 		log.Println(alivePlayers[i])
// 	}
// }

func killPlayer(name string) bool { // takes out player from alive Players list
	for i := 0; i < len(allPlayers); i++ {
		if name == allPlayers[i] {
			// alivePlayers = append(alivePlayers[:i], alivePlayers[i+1:]...) //take out player from list
			allPlayers[i] += " <DEAD> "
			return true
		}
	}
	log.Println("Error! player", name, "not in the alive players list\n")
	return false // player was not found in alive players list
}

//handles information given from public record and handles relevate data
func handleGameString(str string) []byte {
	str = strings.TrimSpace(str)
	commands := strings.Split(str, ";") //split strings by ";" separated values

	finalValue := "" //default output value should be nothingli
	switch {
	case commands[0] == "name":
		log.Println("New player added.")
	case commands[0] == "meta": //a message type for anything else
		if len(commands) == 2 {
			finalValue = fmt.Sprint(commands[1], "\n") //if meta should be read as just a message, then print message
		} else if len(commands) > 2 {
			if commands[1] == "all players" { //if meta is an update of all
				allPlayers = commands[2:] //recreate list of players
				if spectatorMode {
					finalValue = fmt.Sprint("Updated players.\n")
				}
			}
		}
	case commands[0] == "start": //start game
		finalValue = fmt.Sprint("GoWar STARTED!!! \n")
		gameActive = true
	case commands[0] == "stop": //stop game and output scores
		if gameActive {
			finalValue = fmt.Sprint(commands[1]) //print game results given in stop command
			log.Println("\n\n")
			gameActive = false
		}

	case commands[0] == "death": //reports another players death
		if len(commands) < 3 { //error check
			return []byte("error;bad args;death\n")
		}
		finalValue = fmt.Sprint("Player ", commands[1], " has been killed by ", commands[2], "\n") //output who died
		killPlayer(commands[1])
	case commands[0] == "attack": //a player was attacked
		if len(commands) < 4 { //error check
			return []byte("error;bad args;attack\n")
		}
		if spectatorMode { //spectator mode allows players to watch non relevant attacks
			log.Println(commands[1], " was attacked by ", commands[2], " for ", commands[3], "damage.\n")
		} else if commands[1] == myName { //if this player was attacked
			damage, _ := strconv.Atoi(commands[3])
			if myHealth > damage { //if attack but not killed then decrement health and print to user
				myHealth = myHealth - damage
				finalValue = fmt.Sprint("You were attacked by ", commands[2], " for ", commands[3], " damage.\n")
			} else { //if killed notify user, update my killer, report death to system, set health to 0 and enter spectator mode
				finalValue = fmt.Sprint("You were killed by ", commands[2], "!!!\n")
				mykiller = commands[2]
				killPlayer(myName)
				sendDeathMessage = true
				myHealth = 0
				spectatorMode = true
			}
		}
	default:
		// finalValue = "error;unhandled tag;" + commands[0] + "\n"
		finalValue = ""
	}
	return []byte(finalValue) //this value will be given to user
}

//handles messages given from this client, and processes them to be sent to other players
func handleInputString(str string) []byte {
	str = strings.TrimSpace(str)
	commands := strings.Split(str, " ")
	if sendDeathMessage { //if they died then first loop will send death message once
		sendDeathMessage = false
		return []byte(fmt.Sprint("death;", myName, ";", mykiller, "\n"))

	}
	if myHealth == 0 { //check first to see if player is still alive
		log.Print("You are dead :(\n")
		return []byte(fmt.Sprint("\n"))
	}
	finalValue := ""
	switch commands[0] {
	case "help": //print list of usable commands
		log.Print(usageString)
	case "hackyhackhack": //enter custom data without user error thrown, hackers here ya go :)
		finalValue = commands[1] + "\n"
	case "attack": // attack a player, if bad format, or not alive, give error message
		if !gameActive {
			log.Print("Game has not started yet!\n")
		} else if len(commands) < 2 { //bad format error
			log.Print("You have not given enough arguments!\n")
			log.Print("Format: attack [PLAYER NAME]\n")
		} else if isPlayerAlive(commands[1]) { //look to see if that player is alive, if so attack
			atk := rand.Intn(15)
			finalValue = fmt.Sprint("attack;", commands[1], ";", myName, ";", atk, "\n")
			log.Println("Attack successful for ", atk, " damage.")
			playerPoints += atk
		} else { //if player not found then they are dead or were never playing
			log.Print(commands[1], " is dead and gone.")
		}
	case "name": // set name if given proper input
		if len(commands) < 2 {
			log.Println("My name: ", myName) // return user name
			break
		}

		if nameSet { //if name already set, report error since you can only set name once
			log.Print("You can only set your name once!\n")
		} else if len(commands[1]) > 5 { //make sure that name isnt over 5 characters
			log.Print("Your name must be 5 characters or less with no special characters.\n")
		} else {
			// // This is a breaking change becuase we cannot sync the names of all the
			// // players until the game begins

			// found := false
			// for _, v := range allPlayers {
			// 	if v == commands[1] {
			// 		log.Println("Another player has already taken that name!")
			// 		found = true
			// 	}
			// }

			// if found {
			// 	break
			// }

			nameSet = true //set name
			myName = commands[1]
			// allPlayers = append(allPlayers, myName)
			finalValue = fmt.Sprint("name;", commands[1], "\n")
			log.Println("Name has successfully been set.")
		}
	case "list": //lists all players in the server
		if len(allPlayers) == 0 {
			log.Println("Wait till game starts.") //cant know who is in the lobby until game starts
		} else { //list all other players
			log.Print("Players:\n")
			for i := 0; i < len(allPlayers); i++ {
				log.Print("Player ", i, ": ", allPlayers[i])
			}
		}
	case "score": //return players score
		log.Println("My score: ", strconv.Itoa(playerPoints))
	case "spec": //enter spectator mode
		myHealth = 0
		spectatorMode = true
	default: //return error if bad string given
		log.Print("Error: Bad input.\n")
		log.Println(usageString)
		finalValue = fmt.Sprint("error;bad_input_string;", str, "\n")
	}

	return []byte(finalValue) //report user input to system
}

func streamCpy(src io.Reader, dst io.Writer, isOutgoing bool) <-chan int { //reader both from system to client and client to system
	buf := make([]byte, 1024) //message can only be 1024 characters
	sync := make(chan int)

	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				log.Printf("Con from %v closed\n", con.RemoteAddr())
			}
			sync <- 0
		}()

		for { //connect to system read and write functions

			nBytes, err := src.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Read error: %s\n", err)
				}
				break
			}

			// routes the icnoming and outgoing data through the correct
			// functions for interpretation
			var data []byte
			if isOutgoing {
				str := string(buf[0:nBytes])
				data = handleInputString(str)
			} else {
				str := string(buf[0:nBytes])
				data = handleGameString(str)
			}

			// write outgoing data to the destination socket
			_, err = dst.Write(data)

			if err != nil {
				log.Fatalf("Write error: %s\n", err)
			}
		}
	}()

	return sync

}

func HandleCons(con net.Conn) { //connect and disconnect gracefully

	// creates two channels which control the execution of
	// the input and output streams
	stdoutChan := streamCpy(con, os.Stdout, false)
	remoteChan := streamCpy(os.Stdin, con, true)

	// if anything can be pulled from one of the previous channels,
	// the program is over and can now close
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

	rand.Seed(time.Now().UTC().UnixNano()) // set random seed for each client
	var host, port string
	flag.StringVar(&host, "host", "", "Remote host to connect to")
	flag.StringVar(&port, "port", ":8080", "Port of remote host")
	flag.Parse()

	if host == "" {
		flag.Usage()
	}

	StartClient(host, port)

}
