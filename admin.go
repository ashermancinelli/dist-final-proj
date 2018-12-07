package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var (
	players []string
	points	[]int
	alivePlayers     []string
	gameActive       bool = false
)


func givePoints(name string, addpoints int) { //adds player score
	for i := 0; i < len(players); i++ {
		if players[i] == name {
			points[i] = points[i] + addpoints
			return
		}
	}
	log.Print("Error updating ", name, "'s score!")
}
func removePlayer(name string, listName string) { //removes player from  specified list
	if listName == "players" {
		for i := 0; i < len(players); i++ {
			if players[i] == name { //remove player when found
				players = append(players[:i], players[(i+1):]...)
				points = append(points[:i], points[(i+1):]...)
				return
			}
		}
	} else if listName == "alivePlayers" {
		for i := 0; i < len(alivePlayers); i++ {
			if alivePlayers[i] == name { //remove player when found
				alivePlayers = append(alivePlayers[:i], alivePlayers[i+1:]...)
				return
			}
		}
	}
	log.Print("Error deleting player ", name, "!")
}

func handleGameString(str string) []byte { //handles relevant string data from messaging system
	str = strings.TrimSpace(str)
	commands := strings.Split(str, ";") //split strings by ";" separated values

	finalValue := "\n"
	switch {
	case commands[0] == "name": //add new player by name
		if len(commands) < 2 {
			return []byte("error; bad name arg\n")
		}
		players = append(players, commands[1])
		points = append(points, 0)
		alivePlayers = append(alivePlayers, commands[1])
		finalValue = fmt.Sprint("Adding new player: ", commands[1], "\n")
	case commands[0] == "death": //reports another players death
		if len(commands) < 3 { //error check
			return []byte("error;bad args;kill\n")
		}
		finalValue = fmt.Sprint("Player ", commands[2], " has been killed by ", commands[1], "\n") //output who died
		givePoints(commands[3], 10)                                                                 //give 10 point to the killer
		removePlayer(commands[2], "alivePlayers")                                                   //remove the dead person from alive players list
	default:
		finalValue = "error;command_not_implemented\n"
		log.Print("Bad game data: ", str, "\n")
	}
	return []byte(finalValue)
}

func handleInputString(str string) []byte { //valid inputs are start and stop,
	str = strings.TrimSpace(str)
	commands := strings.Split(str, " ")
	finalValue := "\n"
	switch commands[0] {
	case "start": //start the game on admins command
		if gameActive {
			log.Print("game is already active\n")
		} else {
			gameActive = true
			finalValue = fmt.Sprint("start\n")
		}
	case "stop": //quick stop on admins command
		if !gameActive {
			log.Print("game is not yet started\n")
		} else {
			gameActive = false
			finalValue = fmt.Sprint("stop\n")
		}
	case "boot": // boot specific player
		if len(commands) < 2 {
			log.Print("invalid boot command, should be 'boot name'")
			return []byte(finalValue)
		}
		removePlayer(commands[1], "players")
		removePlayer(commands[1], "alivePlayers")
		finalValue = fmt.Sprint(commands[1], " has been booted from game")
	default:
		finalValue = fmt.Sprint(str + "\n")
		log.Print("Custom Admin Input set\n")
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

			if !isOutgoing {
				str := string(buf[0:nBytes])
				data := handleGameString(str)
				_, err = dst.Write(data)
			} else {
				str := string(buf[0:nBytes])
				data := handleInputString(str)
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

func startAdmin(host string, port string) {
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

	startAdmin(host, port)

}
